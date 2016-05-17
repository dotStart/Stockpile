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
import javax.annotation.concurrent.Immutable;
import javax.annotation.concurrent.ThreadSafe;
import java.util.Objects;

/**
 * <strong>Server Information</strong>
 *
 * Provides basic environment information of the server.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Immutable
@ThreadSafe
public class ServerInformation {
        private final Version api;
        private final String name;
        private final String vendor;
        private final String version;

        private ServerInformation() {
                this(new Version(0, Version.State.DEPRECATED));
        }

        public ServerInformation(@Nonnull Version api) {
                this("Stockpile", "Unknown", "Torchmind", api);
        }

        public ServerInformation(String name, String version, String vendor, Version api) {
                this.name = name;
                this.version = version;
                this.vendor = vendor;
                this.api = api;
        }

        public ServerInformation(@Nonnull Package p, @Nonnull Version api) {
                this(p.getImplementationTitle(), p.getImplementationVersion(), p.getImplementationVendor(), api);
        }

        /**
         * Retrieves the API version.
         *
         * @return a version.
         */
        @Nonnull
        public Version getApi() {
                return this.api;
        }

        /**
         * Retrieves the server's implementation name.
         *
         * @return an implementation name.
         */
        @Nonnull
        public String getName() {
                return this.name;
        }

        /**
         * Retrieves the server's vendor.
         *
         * @return a vendor name.
         */
        @Nonnull
        public String getVendor() {
                return this.vendor;
        }

        /**
         * Retrieves the server's implementation version.
         *
         * @return an implementation version.
         */
        @Nonnull
        public String getVersion() {
                return this.version;
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean equals(Object o) {
                if (this == o) return true;
                if (o == null || getClass() != o.getClass()) return false;
                ServerInformation that = (ServerInformation) o;

                return Objects.equals(this.name, that.name) &&
                        Objects.equals(this.version, that.version) &&
                        Objects.equals(this.vendor, that.vendor) &&
                        Objects.equals(this.api, that.api);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public int hashCode() {
                return Objects.hash(this.name, this.version, this.vendor, this.api);
        }
}
