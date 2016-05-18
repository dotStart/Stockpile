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
package com.torchmind.stockpile.server.configuration;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.NestedConfigurationProperty;
import org.springframework.stereotype.Component;

import javax.annotation.Nonnull;

/**
 * <strong>Cache Configuration</strong>
 *
 * Represents an externalized cache configuration.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Component
@ConfigurationProperties(prefix = "cache")
public class CacheConfiguration {
        private Aggressiveness aggressiveness;
        @NestedConfigurationProperty
        private TimeToLiveConfiguration ttl = new TimeToLiveConfiguration();

        @Nonnull
        public Aggressiveness getAggressiveness() {
                return this.aggressiveness;
        }

        @Nonnull
        public TimeToLiveConfiguration getTtl() {
                return this.ttl;
        }

        public void setAggressiveness(@Nonnull Aggressiveness aggressiveness) {
                this.aggressiveness = aggressiveness;
        }

        public void setTtl(@Nonnull TimeToLiveConfiguration ttl) {
                this.ttl = ttl;
        }

        /**
         * <strong>Cache Aggressiveness</strong>
         *
         * Provides a list of general behavior presets.
         */
        public enum Aggressiveness {

                /**
                 * <strong>Low Aggressiveness</strong>
                 *
                 * Switches the cache to write-only mode unless the backing API returns an error.
                 */
                LOW,

                /**
                 * <strong>Moderate Aggressiveness</strong>
                 *
                 * Enables user-customizable caching using the TTL options.
                 */
                MODERATE,

                /**
                 * <strong>High Aggressiveness</strong>
                 *
                 * Ensures all objects are cached for the longest possible time. Only push updates will alter the set of
                 * local records.
                 */
                HIGH,

                /**
                 * <strong>Unknown</strong>
                 *
                 * This is a placeholder annotation for cases in which unknown or default values are required.
                 */
                UNKNOWN
        }

        /**
         * Represents an externalized TTL configuration.
         */
        public static class TimeToLiveConfiguration {
                private long name = 7200;
                private long profile = -1;
                private long property = 3600;

                public long getName() {
                        return this.name;
                }

                public long getProfile() {
                        return this.profile;
                }

                public long getProperty() {
                        return this.property;
                }

                public void setName(long name) {
                        this.name = name;
                }

                public void setProfile(long profile) {
                        this.profile = profile;
                }

                public void setProperty(long property) {
                        this.property = property;
                }
        }
}
