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
package com.torchmind.stockpile.server.service.api;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectReader;
import com.torchmind.stockpile.server.configuration.CacheConfiguration;
import com.torchmind.stockpile.server.entity.DisplayName;
import com.torchmind.stockpile.server.entity.Profile;
import com.torchmind.stockpile.server.entity.ProfileProperty;
import com.torchmind.stockpile.server.entity.repository.DisplayNameRepository;
import com.torchmind.stockpile.server.entity.repository.ProfilePropertyRepository;
import com.torchmind.stockpile.server.entity.repository.ProfileRepository;
import com.torchmind.stockpile.server.error.NoSuchProfileException;
import com.torchmind.stockpile.server.error.ServiceException;
import com.torchmind.stockpile.server.error.TooManyRequestsException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.annotation.Nonnull;
import javax.annotation.Nullable;
import javax.annotation.concurrent.ThreadSafe;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

/**
 * <strong>Profile Service</strong>
 *
 * Provides simplified methods of retrieving a specific profile
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Service
@ThreadSafe
public class ProfileService {
        public static final String NAME_URL_TEMPLATE = "https://api.mojang.com/users/profiles/minecraft/%s";
        public static final String PROFILE_URL_TEMPLATE = "https://sessionserver.mojang.com/session/minecraft/profile/%s?unsigned=false";
        public static final ObjectReader reader;

        private final CacheConfiguration cacheConfiguration;
        private final DisplayNameRepository displayNameRepository;
        private final ProfilePropertyRepository profilePropertyRepository;
        private final ProfileRepository profileRepository;

        static {
                ObjectMapper mapper = new ObjectMapper();
                mapper.findAndRegisterModules();
                reader = mapper.reader();
        }

        @Autowired
        public ProfileService(@Nonnull CacheConfiguration cacheConfiguration, @Nonnull ProfileRepository profileRepository, @Nonnull DisplayNameRepository displayNameRepository, @Nonnull ProfilePropertyRepository profilePropertyRepository) {
                this.cacheConfiguration = cacheConfiguration;
                this.profileRepository = profileRepository;
                this.displayNameRepository = displayNameRepository;
                this.profilePropertyRepository = profilePropertyRepository;
        }

        /**
         * Fetches an entire profile directly from Mojang and adds it to the local cache.
         *
         * @param identifier a profile identifier.
         * @return a profile.
         *
         * @throws IOException              when an error occurs while contacting Mojang.
         * @throws TooManyRequestsException when the rate limit is exceeded.
         */
        @Nonnull
        private Profile fetch(@Nonnull UUID identifier) throws IOException {
                URL profileUrl = new URL(String.format(PROFILE_URL_TEMPLATE, (new MojangUUID(identifier)).toString()));
                HttpURLConnection connection = (HttpURLConnection) profileUrl.openConnection();

                if (connection.getResponseCode() == 429) {
                        throw new TooManyRequestsException("Rate limit exceeded");
                }

                try (InputStream inputStream = connection.getInputStream()) {
                        JsonNode node = reader.readTree(inputStream);

                        // try to fetch a profile or create a new one if none is stored within the database
                        final Profile profile;
                        {
                                Profile prof = this.profileRepository.findOne(identifier);

                                if (prof == null) {
                                        profile = new Profile(identifier);
                                } else {
                                        profile = prof;
                                        profile.setLastSeen(Instant.now());
                                }
                        }
                        this.profileRepository.save(profile);

                        // try to fetch a display name for the current profile and create a new one if none is stored in
                        // the database
                        DisplayName displayName = this.displayNameRepository.findOneByNameAndProfile(node.get("name").asText(), profile).orElseGet(() -> new DisplayName(node.get("name").asText(), Instant.now(), profile));

                        displayName.setLastSeen(Instant.now());
                        this.displayNameRepository.save(displayName);

                        // iterate over all properties and create/update their respective values
                        node.get("properties").forEach((p) -> {
                                String name = p.get("name").asText();

                                ProfileProperty property = this.profilePropertyRepository.findByNameAndProfile(name, profile).orElseGet(() -> new ProfileProperty(profile, name, p.get("value").asText(), (p.has("signature") ? p.get("signature").asText() : null)));

                                property.setLastSeen(Instant.now());
                                this.profilePropertyRepository.save(property);
                                profile.addProperty(property);
                        });

                        return profile;
                }
        }

        /**
         * Fetches a display name's associated identifier directly from Mojang.
         *
         * @param name a display name.
         * @return an identifier.
         *
         * @throws IOException              when an error occurs while contacting Mojang.
         * @throws TooManyRequestsException when the rate limit is exceeded.
         */
        @Nonnull
        private UUID fetchIdentifier(@Nonnull String name) throws IOException {
                URL identifierUrl = new URL(String.format(NAME_URL_TEMPLATE, name));
                HttpURLConnection connection = (HttpURLConnection) identifierUrl.openConnection();

                if (connection.getResponseCode() == 429) {
                        throw new TooManyRequestsException("Rate limit exceeded");
                }

                try (InputStream inputStream = connection.getInputStream()) {
                        JsonNode node = reader.readTree(inputStream);
                        UUID identifier = (new MojangUUID(node.get("id").asText())).toUUID();

                        // fetch existing profile or create an entirely new one
                        final Profile profile;
                        {
                                Profile prof = this.profileRepository.findOne(identifier);

                                if (prof == null) {
                                        profile = new Profile(identifier);
                                } else {
                                        profile = prof;
                                        profile.setLastSeen(Instant.now());
                                }
                        }
                        this.profileRepository.save(profile);

                        // fetch an existing display name or create an entirely new one
                        final String currentName = node.get("name").asText();
                        final DisplayName displayName = this.displayNameRepository.findOneByNameAndProfile(currentName, profile).orElseGet(() -> new DisplayName(currentName, Instant.now(), profile));
                        displayName.setLastSeen(Instant.now());
                        this.displayNameRepository.save(displayName);
                        profile.addName(displayName);

                        return identifier;
                }
        }

        /**
         * Fetches a cached version of a display name identifier from the database backend.
         *
         * @param name a display name.
         * @return an identifier or null.
         */
        @Nullable
        private UUID fetchIdentifierLocal(@Nonnull String name) {
                return this.displayNameRepository.findOneByName(name).map((d) -> d.getProfile().getIdentifier()).orElse(null);
        }

        /**
         * Fetches a cached version of a profile from the database backend.
         *
         * @param identifier an identifier.
         * @return a profile or null.
         */
        @Nullable
        private Profile fetchLocal(@Nonnull UUID identifier) {
                return this.profileRepository.findOne(identifier);
        }

        /**
         * Finds a profile based on its current display name in the local cache or directly in Mojang's database based
         * on the current cache aggressiveness.
         *
         * @param name a display name.
         * @return a profile.
         *
         * @throws NoSuchProfileException when the profile was not found.
         * @throws ServiceException       when the profile could not be accessed.
         */
        @Nonnull
        public Profile find(@Nonnull String name) {
                return this.get(this.findIdentifier(name));
        }

        /**
         * Finds a profile identifier based on its current display name in the local cache or directly in Mojang's
         * database based on the current cache aggressiveness.
         *
         * @param name a display name.
         * @return a profile.
         */
        @Nonnull
        public UUID findIdentifier(@Nonnull String name) {
                if (this.cacheConfiguration.getAggressiveness() == CacheConfiguration.Aggressiveness.LOW) {
                        UUID identifier;

                        try {
                                identifier = this.fetchIdentifier(name);
                        } catch (TooManyRequestsException | IOException ex) {
                                identifier = this.fetchIdentifierLocal(name);
                        }

                        if (identifier == null) {
                                throw new NoSuchProfileException(name);
                        }

                        return identifier;
                }

                UUID identifier = this.fetchIdentifierLocal(name);

                if (identifier == null) {
                        try {
                                identifier = this.fetchIdentifier(name);
                        } catch (FileNotFoundException ex) {
                                throw new NoSuchProfileException(name);
                        } catch (TooManyRequestsException | IOException ex) {
                                throw new ServiceException("Could not poll identifier from upstream: " + ex.getMessage(), ex);
                        }
                }

                return identifier;
        }

        /**
         * Retrieves a profile from the cache or pulls a fresh copy directly from Mojang based on the current cache
         * aggressiveness.
         *
         * @param identifier an identifier.
         * @return a profile.
         *
         * @throws NoSuchProfileException when the profile was not found.
         * @throws ServiceException       when the profile could not be accessed.
         */
        @Nonnull
        public Profile get(@Nonnull UUID identifier) {
                if (this.cacheConfiguration.getAggressiveness() == CacheConfiguration.Aggressiveness.LOW) {
                        try {
                                this.fetch(identifier);
                        } catch (FileNotFoundException ex) {
                                throw new NoSuchProfileException(identifier);
                        } catch (TooManyRequestsException | IOException ex) {
                                return Optional.ofNullable(this.fetchLocal(identifier)).orElseThrow(() -> new ServiceException("Cannot poll nor find cached profile \"" + identifier + "\": " + ex.getMessage(), ex));
                        }
                }

                Profile profile = this.fetchLocal(identifier);

                if (profile == null || profile.getNames().size() == 0 || profile.getProperties().size() == 0) {
                        try {
                                return this.fetch(identifier);
                        } catch (TooManyRequestsException ex) {
                                throw ex;
                        } catch (IOException ex) {
                                if (profile != null) {
                                        return profile;
                                }

                                throw new ServiceException("Could not poll profile from upstream: " + ex.getMessage(), ex);
                        }
                }

                return profile;
        }

        /**
         * Purges an entire profile cache and all of its associated data.
         *
         * @param identifier an identifier.
         */
        public void purge(@Nonnull UUID identifier) {
                this.profileRepository.delete(identifier);
        }

        /**
         * Purges an entire profile cache and all of its associated data.
         *
         * @param name a display name.
         */
        public void purge(@Nonnull String name) {
                this.displayNameRepository.findOneByName(name).ifPresent((n) -> this.profileRepository.delete(n.getProfile()));
        }
}
