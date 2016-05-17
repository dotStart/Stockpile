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
package com.torchmind.stockpile.data;

import javax.annotation.Nonnull;
import java.util.Objects;

/**
 * <strong>Version</strong>
 *
 * Represents the Stockpile API version as well as its state.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public class Version {
        private final State state;
        private final int version;

        public Version(int version, State state) {
                this.version = version;
                this.state = state;
        }

        /**
         * Retrieves the API state.
         *
         * @return a state.
         *
         * @see State for more information on API version states.
         */
        @Nonnull
        public State getState() {
                return state;
        }

        /**
         * Retrieves the numeric version representation.
         *
         * The numeric version is bumped whenever a new feature is added or a backwards incompatible change is made.
         * A server may additionally support multiple versions.
         *
         * @return a numeric version.
         */
        public int getVersion() {
                return version;
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean equals(Object o) {
                if (this == o) return true;
                if (o == null || this.getClass() != o.getClass()) return false;
                Version version1 = (Version) o;

                return this.version == version1.version && this.state == version1.state;
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public int hashCode() {
                return Objects.hash(this.state, this.version);
        }

        /**
         * <strong>API Version State</strong>
         *
         * Provides a list of valid API version states.
         */
        public enum State {

                /**
                 * <strong>Development</strong>
                 *
                 * Signifies that the API is still in development and thus considered unstable (e.g. it may change at
                 * any time and is more likely to produce errors or false results in edge cases).
                 */
                DEVELOPMENT,

                /**
                 * <strong>Stable</strong>
                 *
                 * Signifies that the API has been finalized and is thus considered stable (e.g. it will not be changed
                 * anymore and is unlikely to produce errors or false results in edge cases).
                 */
                STABLE,

                /**
                 * <strong>Deprecated</strong>
                 *
                 * Signifies that the API has been deprecated and thus clients should refrain from using it.
                 *
                 * Deprecated API versions may be removed in future versions of the Stockpile server implementation and
                 * may thus cause clients to not be able to connect to the server when an update is installed.
                 *
                 * Generally client implementors are expected to warn users of this circumstance when initializing to
                 * ensure updates are installed in a timely manner.
                 */
                DEPRECATED
        }
}
