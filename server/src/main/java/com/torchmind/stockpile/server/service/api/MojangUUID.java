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

import javax.annotation.Nonnull;
import java.util.Objects;
import java.util.UUID;

/**
 * <strong>Mojang UUID</strong>
 *
 * Represents a Mojang encoded UUID (a UUID which does not include dashes).
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public class MojangUUID {
        private final String encoded;

        public MojangUUID(@Nonnull String encoded) {
                this.encoded = encoded;
        }

        public MojangUUID(@Nonnull UUID uuid) {
                this.encoded = uuid.toString().replace("-", "");
        }

        /**
         * Converts a Mojang UUID into a regular sane UUID.
         *
         * @return a UUID.
         */
        @Nonnull
        public UUID toUUID() {
                String group1 = this.encoded.substring(0, 8);
                String group2 = this.encoded.substring(8, 12);
                String group3 = this.encoded.substring(12, 16);
                String group4 = this.encoded.substring(16, 20);
                String group5 = this.encoded.substring(20);

                return UUID.fromString(group1 + '-' + group2 + '-' + group3 + '-' + group4 + '-' + group5);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean equals(Object o) {
                if (this == o) return true;
                if (o == null || this.getClass() != o.getClass()) return false;

                MojangUUID that = (MojangUUID) o;
                return Objects.equals(this.encoded, that.encoded);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public int hashCode() {
                return Objects.hash(this.encoded);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public String toString() {
                return this.encoded;
        }
}
