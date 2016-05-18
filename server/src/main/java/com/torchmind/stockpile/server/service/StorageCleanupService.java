/*
 * Copyright 2016 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.torchmind.stockpile.server.service;

import com.torchmind.stockpile.server.configuration.CacheConfiguration;
import com.torchmind.stockpile.server.configuration.condition.ConditionalOnCacheAggressiveness;
import com.torchmind.stockpile.server.entity.DisplayName;
import com.torchmind.stockpile.server.entity.Profile;
import com.torchmind.stockpile.server.entity.ProfileProperty;
import com.torchmind.stockpile.server.entity.repository.DisplayNameRepository;
import com.torchmind.stockpile.server.entity.repository.ProfilePropertyRepository;
import com.torchmind.stockpile.server.entity.repository.ProfileRepository;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import javax.annotation.Nonnull;
import javax.annotation.PostConstruct;
import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * <strong>Storage Cleanup Service</strong>
 *
 * Provides methods for cleaning up old cached values.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public interface StorageCleanupService {

        /**
         * <strong>Abstract Cleanup Service</strong>
         *
         * Provides a base implementation for cleanup services.
         */
        abstract class AbstractCleanupService implements StorageCleanupService {
                protected final Logger logger = LogManager.getFormatterLogger(this.getClass());

                /**
                 * Executes a database cleanup.
                 */
                protected abstract void cleanup();

                /**
                 * Provides a common scheduler registration.
                 *
                 * @see #cleanup()
                 */
                @Scheduled(cron = "0 0/5 * * * ?")
                public final void doCleanup() {
                        this.cleanup();
                }

                /**
                 * Ensures database cleanups are executed at least once during application startup.
                 */
                @PostConstruct
                public final void doSetup() {
                        this.setup();

                        this.logger.info("Commencing initial database cleanup.");
                        this.cleanup();
                }

                /**
                 * Handles the service initialization.
                 */
                protected void setup() {
                }
        }

        /**
         * <strong>High Aggressiveness Cleanup Service</strong>
         *
         * Cleans out all objects that reached the maximum possible TTL.
         */
        @Service
        @ConditionalOnCacheAggressiveness(CacheConfiguration.Aggressiveness.HIGH)
        class HighAggressivenessCleanupService extends AbstractCleanupService {
                public static final Duration NAME_EXPIRATION_DURATION = Duration.ofDays(30).minusMinutes(6); // 6 Minutes less than actual timeout to work around issues with the scheduler
                private final DisplayNameRepository displayNameRepository;

                @Autowired
                public HighAggressivenessCleanupService(@Nonnull DisplayNameRepository displayNameRepository) {
                        this.displayNameRepository = displayNameRepository;
                }

                /**
                 * {@inheritDoc}
                 */
                @Override
                public void cleanup() {
                        long count = 0;
                        logger.info("Cleaning up cached values using high aggressiveness ...");
                        {
                                try (Stream<DisplayName> stream = this.displayNameRepository.findByLastSeenLessThanOrderByLastSeenDesc(Instant.now().minus(NAME_EXPIRATION_DURATION))) {
                                        List<DisplayName> names = stream.collect(Collectors.toList());

                                        count = names.size();
                                        this.displayNameRepository.delete(names);
                                }
                        }
                        logger.info("Cleanup successful: %d values have been removed.", count);
                }

        }

        /**
         * <strong>Low Aggressiveness Cleanup Service</strong>
         *
         * Cleans out all objects that reached their maximum TTL.
         */
        @Service
        @ConditionalOnCacheAggressiveness(CacheConfiguration.Aggressiveness.LOW)
        class LowAggressivenessCleanupService extends HighAggressivenessCleanupService {

                @Autowired
                public LowAggressivenessCleanupService(@Nonnull DisplayNameRepository displayNameRepository) {
                        super(displayNameRepository);
                }
        }

        /**
         * <strong>Moderate Aggressiveness Cleanup Service</strong>
         *
         * Cleans out all objects that reached the configured TTL.
         */
        @Service
        @ConditionalOnCacheAggressiveness(CacheConfiguration.Aggressiveness.MODERATE)
        class ModerateAggressivenessCleanupService extends AbstractCleanupService {
                private final CacheConfiguration cacheConfiguration;
                private final DisplayNameRepository displayNameRepository;
                private final ProfilePropertyRepository profilePropertyRepository;
                private final ProfileRepository profileRepository;

                @Autowired
                public ModerateAggressivenessCleanupService(@Nonnull CacheConfiguration cacheConfiguration, @Nonnull DisplayNameRepository displayNameRepository, @Nonnull ProfilePropertyRepository profilePropertyRepository, @Nonnull ProfileRepository profileRepository) {
                        this.cacheConfiguration = cacheConfiguration;
                        this.displayNameRepository = displayNameRepository;
                        this.profilePropertyRepository = profilePropertyRepository;
                        this.profileRepository = profileRepository;
                }

                /**
                 * {@inheritDoc}
                 */
                @Override
                public void cleanup() {
                        long count = 0;
                        logger.info("Cleaning up cached values using moderate aggressiveness ...");
                        {
                                if (this.cacheConfiguration.getTtl().getProfile() > 0) {
                                        count += this.cleanupProfiles();
                                }

                                if (this.cacheConfiguration.getTtl().getName() > 0) {
                                        count += this.cleanupDisplayNames();
                                }

                                if (this.cacheConfiguration.getTtl().getProperty() > 0) {
                                        count += this.cleanupProperties();
                                }
                        }
                        logger.info("Cleanup successful: %d values have been removed.", count);
                }

                /**
                 * Removes all profiles which are considered expired based on the user's TTL configuration.
                 *
                 * @return the amount of deleted display names.
                 */
                private long cleanupDisplayNames() {
                        Instant expirationTimestamp = Instant.now().minusSeconds(this.cacheConfiguration.getTtl().getName());

                        try (Stream<DisplayName> stream = this.displayNameRepository.findByLastSeenLessThanOrderByLastSeenDesc(expirationTimestamp)) {
                                List<DisplayName> names = stream.collect(Collectors.toList());
                                long count = names.size();

                                this.displayNameRepository.delete(names);
                                return count;
                        }
                }

                /**
                 * Removes all profiles which are considered expired based on the user's TTL configuration.
                 *
                 * @return the amount of deleted profiles.
                 */
                private long cleanupProfiles() {
                        Instant expirationTimestamp = Instant.now().minusSeconds(this.cacheConfiguration.getTtl().getProfile());

                        try (Stream<Profile> stream = this.profileRepository.findByLastSeenLessThanOrderByLastSeenDesc(expirationTimestamp)) {
                                List<Profile> profiles = stream.collect(Collectors.toList());
                                long count = profiles.size();

                                this.profileRepository.delete(profiles);
                                return count;
                        }
                }

                /**
                 * Removes all profiles which are considered expired based on the user's TTL configuration.
                 *
                 * @return the amount of deleted properties.
                 */
                private long cleanupProperties() {
                        Instant expirationTimestamp = Instant.now().minusSeconds(this.cacheConfiguration.getTtl().getProperty());

                        try (Stream<ProfileProperty> stream = this.profilePropertyRepository.findByLastSeenLessThanOrderByLastSeenDesc(expirationTimestamp)) {
                                List<ProfileProperty> properties = stream.collect(Collectors.toList());
                                long count = properties.size();

                                this.profilePropertyRepository.delete(properties);
                                return count;
                        }
                }

                /**
                 * {@inheritDoc}
                 */
                @Override
                protected void setup() {
                        if (this.cacheConfiguration.getTtl().getName() > Duration.ofDays(30).getSeconds()) {
                                logger.warn("+=============================+");
                                logger.warn("| DANGEROUS TTL CONFIGURATION |");
                                logger.warn("+-----------------------------+");
                                logger.warn("| The server was configured   |");
                                logger.warn("| to cache display names for  |");
                                logger.warn("| more than 30 days!          |");
                                logger.warn("|                             |");
                                logger.warn("| This WILL cause issues when |");
                                logger.warn("| searching for profiles or   |");
                                logger.warn("| UUIDs based on names!       |");
                                logger.warn("+=============================+");
                        }
                }
        }
}
