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

import javax.annotation.Nonnull;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.ManyToOne;
import javax.persistence.Table;

/**
 * <strong>Profile Property</strong>
 *
 * Represents a property which has been assigned to a profile.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Entity
@Table(name = "profile_property")
public class ProfileProperty extends BaseEntity {
        @Column(nullable = false, updatable = false)
        private final String name;
        @Column(nullable = false)
        private String value;
        @Column
        private String signature;
        @ManyToOne(optional = false)
        private final Profile profile;

        private ProfileProperty() {
                this.name = null;
                this.profile = null;
        }

        public ProfileProperty(@Nonnull Profile profile, @Nonnull String name, @Nonnull String value, @Nonnull String signature) {
                this.profile = profile;

                this.name = name;
                this.value = value;
                this.signature = signature;
        }

        @Nonnull
        public String getName() {
                return this.name;
        }

        @Nonnull
        public String getValue() {
                return this.value;
        }

        public void setValue(@Nonnull String value) {
                this.value = value;
        }

        @Nonnull
        public String getSignature() {
                return this.signature;
        }

        public void setSignature(@Nonnull String signature) {
                this.signature = signature;
        }

        @Nonnull
        public Profile getProfile() {
                return this.profile;
        }
}
