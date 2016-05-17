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

/**
 * <strong>Blacklist Result</strong>
 *
 * Represents the answer to a blacklist check.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Immutable
@ThreadSafe
public class BlacklistResult {
        private final boolean blacklisted;
        private final String hostname;

        private BlacklistResult() {
                this("*.example.org", false);
        }

        public BlacklistResult(@Nonnull String hostname, boolean blacklisted) {
                this.hostname = hostname;
                this.blacklisted = blacklisted;
        }

        /**
         * Retrieves the blacklisted hostname as stored in the database.
         *
         * @return a hostname (may include wildcards).
         */
        @Nonnull
        public String getHostname() {
                return hostname;
        }

        /**
         * Checks whether the host is blacklisted.
         *
         * @return true if blacklisted, false otherwise.
         */
        public boolean isBlacklisted() {
                return blacklisted;
        }
}
