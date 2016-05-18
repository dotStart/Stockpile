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
package com.torchmind.stockpile.data.v1;

import javax.annotation.Nonnull;
import javax.annotation.Nullable;
import javax.annotation.concurrent.Immutable;
import javax.annotation.concurrent.ThreadSafe;
import java.time.Instant;
import java.util.*;

/**
 * <strong>Player Profile</strong>
 *
 * Represents a (possibly cached) player profile.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Immutable
@ThreadSafe
public class PlayerProfile {
        private final Instant cacheTimestamp;
        private final UUID identifier;
        private final String name;
        private final List<Property> properties;

        private PlayerProfile() {
                this(new UUID(0, 0), "Invalid", Collections.emptyList(), null);
        }

        public PlayerProfile(@Nonnull UUID identifier, @Nonnull String name, @Nonnull List<Property> properties, @Nonnull Instant cacheTimestamp) {
                this.identifier = identifier;
                this.name = name;
                this.properties = properties;
                this.cacheTimestamp = cacheTimestamp;
        }

        /**
         * Searches a property within the profile record.
         *
         * @param name a profile name.
         * @return a property or, if no property with the supplied name was found within the record, an empty optional.
         */
        @Nonnull
        public Optional<Property> findProperty(@Nonnull String name) {
                return this.properties.stream().filter((p) -> p.getName().equals(name)).findAny();
        }

        /**
         * Retrieves the timestamp this profile was last updated at.
         * <strong>Note:</strong> This timestamp will never be more than 30 days in the past.
         *
         * @return a timestamp.
         */
        @Nonnull
        public Instant getCacheTimestamp() {
                return this.cacheTimestamp;
        }

        /**
         * Retrieves the profile's globally unique identifier.
         *
         * @return an identifier.
         */
        @Nonnull
        public UUID getIdentifier() {
                return this.identifier;
        }

        /**
         * Retrieves the profile's (cached) display name.
         *
         * @return a display name.
         */
        @Nonnull
        public String getName() {
                return this.name;
        }

        /**
         * Retrieves a list of properties which have been assigned to the profile.
         *
         * @return an immutable set of properties.
         */
        @Nonnull
        public List<Property> getProperties() {
                return Collections.unmodifiableList(this.properties);
        }

        /**
         * <strong>Player Profile Property</strong>
         *
         * Represents a single property within a player profile such as a skin texture.
         */
        @Immutable
        @ThreadSafe
        public static class Property {
                private final String name;
                private final String signature;
                private final String value;

                private Property() {
                        this("Invalid", "", null);
                }

                public Property(@Nonnull String name, @Nonnull String value, @Nullable String signature) {
                        this.name = name;
                        this.value = value;
                        this.signature = signature;
                }

                /**
                 * Retrieves the property's name.
                 *
                 * @return a name.
                 */
                @Nonnull
                public String getName() {
                        return this.name;
                }

                /**
                 * Retrieves the property's signature (if supplied).
                 *
                 * @return a signature or, if no value was cached, an empty optional.
                 */
                @Nullable
                public String getSignature() {
                        return this.signature;
                }

                /**
                 * Retrieves the property's value.
                 *
                 * @return a value.
                 */
                @Nonnull
                public String getValue() {
                        return this.value;
                }

                /**
                 * {@inheritDoc}
                 */
                @Override
                public boolean equals(Object o) {
                        if (this == o) return true;
                        if (o == null || this.getClass() != o.getClass()) return false;

                        Property property = (Property) o;
                        return Objects.equals(this.name, property.name) &&
                                Objects.equals(this.signature, property.signature) &&
                                Objects.equals(this.value, property.value);
                }

                /**
                 * {@inheritDoc}
                 */
                @Override
                public int hashCode() {
                        return Objects.hash(this.name, this.signature, this.value);
                }
        }
}
