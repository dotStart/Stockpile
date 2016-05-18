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
package com.torchmind.stockpile.server.entity;

import com.torchmind.stockpile.data.v1.PlayerProfile;

import javax.annotation.Nonnull;
import javax.persistence.*;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CopyOnWriteArraySet;

/**
 * <strong>Profile</strong>
 *
 * Represents a cached Minecraft profile.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Entity
@Table(name = "profile")
public class Profile {
        @Transient
        private final transient boolean cached;
        @Id
        private final UUID identifier;
        @Column
        private Instant lastSeen;
        @OneToMany(orphanRemoval = true, mappedBy = "profile", fetch = FetchType.EAGER, cascade = CascadeType.REFRESH)
        private final Set<DisplayName> names;
        @OneToMany(orphanRemoval = true, mappedBy = "profile", fetch = FetchType.EAGER, cascade = CascadeType.REFRESH)
        private final Set<ProfileProperty> properties;

        private Profile() {
                this.identifier = null;
                this.names = null;
                this.properties = null;
                this.lastSeen = Instant.now();

                this.cached = true;
        }

        public Profile(@Nonnull UUID identifier) {
                this.identifier = identifier;
                this.names = new CopyOnWriteArraySet<>();
                this.properties = new CopyOnWriteArraySet<>();
                this.lastSeen = Instant.now();

                this.cached = false;
        }

        public void addName(@Nonnull DisplayName name) {
                this.names.add(name);
        }

        public void addProperty(@Nonnull ProfileProperty property) {
                this.properties.add(property);
        }

        @Nonnull
        public UUID getIdentifier() {
                return identifier;
        }

        @Nonnull
        public Instant getLastSeen() {
                return lastSeen;
        }

        /**
         * Retrieves the latest profile name.
         *
         * @return a name or, if no display name was found, an empty optional.
         */
        @Nonnull
        public Optional<DisplayName> getLatestName() {
                return this.getNames().stream().sorted((n1, n2) -> (int) Math.min(1, Math.max(-1, (n2.getLastSeen().toEpochMilli() - n1.getLastSeen().toEpochMilli())))).findFirst();
        }

        /**
         * Retrieves a list of associated names.
         *
         * @return a list of names.
         */
        @Nonnull
        public Set<DisplayName> getNames() {
                return this.names;
        }

        /**
         * Retrieves a list of associated properties.
         *
         * @return a list of properties.
         */
        @Nonnull
        public Set<ProfileProperty> getProperties() {
                return this.properties;
        }

        /**
         * Checks whether the profile has just been created as a result of the request or whether it was received
         * from the storage backend.
         *
         * @return true if cached, false otherwise.
         */
        public boolean isCached() {
                return this.cached;
        }

        public void setLastSeen(@Nonnull Instant lastSeen) {
                this.lastSeen = lastSeen;
        }

        /**
         * Converts a stored profile into a REST compatible resource.
         *
         * @return a REST profile.
         */
        @Nonnull
        public PlayerProfile toRestRepresentation() {
                List<PlayerProfile.Property> properties = new ArrayList<>();
                this.getProperties().forEach((p) -> properties.add(p.toRestRepresentation()));

                return new PlayerProfile(this.getIdentifier(), this.getLatestName().map((n) -> n.getName()).orElse(null), properties, this.getLastSeen());
        }
}
