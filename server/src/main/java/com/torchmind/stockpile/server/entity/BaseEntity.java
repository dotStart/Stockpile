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

import org.hibernate.annotations.GenericGenerator;

import javax.annotation.Nonnull;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.MappedSuperclass;
import java.util.Objects;
import java.util.UUID;

/**
 * <strong>Base Entity</strong>
 *
 * Provides a base type for entities.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@MappedSuperclass
public class BaseEntity {
        @Id
        @GeneratedValue(generator = "uuid")
        @GenericGenerator(name = "uuid", strategy = "uuid2")
        private final UUID identifier;

        protected BaseEntity() {
                this.identifier = null;
        }

        protected BaseEntity(@Nonnull UUID identifier) {
                this.identifier = identifier;
        }

        /**
         * Retrieves the entity's globally unique identifier.
         *
         * @return an identifier.
         */
        public UUID getIdentifier() {
                return identifier;
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean equals(Object o) {
                if (this == o) return true;
                if (!(o instanceof BaseEntity)) return false;

                BaseEntity that = (BaseEntity) o;
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
