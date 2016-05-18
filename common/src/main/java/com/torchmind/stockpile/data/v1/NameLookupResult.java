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
import java.util.Objects;
import java.util.UUID;

/**
 * <strong>Name Lookup Result</strong>
 *
 * Represents a name lookup result which contains the latest known identifier.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public class NameLookupResult {
        private final UUID identifier;

        @SuppressWarnings("ConstantConditions")
        private NameLookupResult() {
                this(null);
        }

        public NameLookupResult(@Nonnull UUID identifier) {
                this.identifier = identifier;
        }

        /**
         * Retrieves the last known associated identifier.
         * @return an identifier.
         */
        @Nonnull
        public UUID getIdentifier() {
                return this.identifier;
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean equals(Object o) {
                if (this == o) return true;
                if (o == null || this.getClass() != o.getClass()) return false;

                NameLookupResult that = (NameLookupResult) o;
                return Objects.equals(this.identifier, that.identifier);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public int hashCode() {
                return Objects.hash(this.identifier);
        }
}
