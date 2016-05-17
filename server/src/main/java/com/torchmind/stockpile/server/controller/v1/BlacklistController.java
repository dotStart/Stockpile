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
package com.torchmind.stockpile.server.controller.v1;

import com.google.common.base.Joiner;
import com.google.common.base.Splitter;
import com.google.common.hash.Hashing;
import com.torchmind.stockpile.data.v1.BlacklistResult;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Nonnull;
import javax.annotation.PostConstruct;
import javax.annotation.concurrent.ThreadSafe;
import java.io.*;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.regex.Pattern;

/**
 * <strong>Blacklist Controller</strong>
 *
 * Provides a controller which allows directly checking against a cached version of the blocked server list.
 * Note: This controller is hardcoded to update its caches once an hour and is not affected by any cache aggressiveness
 * settings.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@ThreadSafe
@RestController
@RequestMapping("/v1/blacklist")
public class BlacklistController {
        public static final String BLACKLIST_URL = "https://sessionserver.mojang.com/blockedservers";
        public static final Pattern IP_ADDRESS_PATTERN = Pattern.compile("^((?:([0-9]{1,2})|(?:([0-1][0-9]{2})|(2(?:([0-4][0-9])|(5[0-5])))))\\.){3}(?:(?:([0-1][0-9]{2})|(2(?:([0-4][0-9])|(5[0-5]))))|([0-9]{1,2}))$");
        public static final Logger logger = LogManager.getFormatterLogger(BlacklistController.class);

        private final List<String> hashes = new CopyOnWriteArrayList<>();

        /**
         * <code>POST /v1/blacklist/</code>
         *
         * Checks any hostname supplied as a form parameter in the post body against the server blacklist.
         *
         * @param hostname a hostname to check against.
         * @return a blacklist result.
         */
        @Nonnull
        @RequestMapping(params = "hostname", method = RequestMethod.POST)
        public BlacklistResult check(@Nonnull @RequestParam("hostname") String hostname) {
                // before checking for wildcards check for exact matches of the hostname
                hostname = hostname.toLowerCase();
                String hash = Hashing.sha1().hashString(hostname, StandardCharsets.ISO_8859_1).toString();

                if (this.hashes.contains(hash)) {
                        return new BlacklistResult(hostname, true);
                }

                if (IP_ADDRESS_PATTERN.matcher(hostname).matches()) {
                        return this.checkAddress(hostname);
                }

                return this.checkHostname(hostname);
        }

        /**
         * Checks an IP address against the blacklist.
         *
         * @param address an address.
         * @return a blacklist result.
         */
        @Nonnull
        private BlacklistResult checkAddress(@Nonnull String address) {
                List<String> addressParts = Splitter.on('.').splitToList(address);

                for (int i = (addressParts.size() - 1); i >= 1; --i) {
                        String currentAddress = Joiner.on('.').join(addressParts.subList(0, i)) + ".*";
                        String hash = Hashing.sha1().hashString(currentAddress, StandardCharsets.ISO_8859_1).toString();

                        if (this.hashes.contains(hash)) {
                                return new BlacklistResult(currentAddress, true);
                        }
                }

                return new BlacklistResult(address, false);
        }

        /**
         * <code>POST /v1/blacklist/</code>
         *
         * Checks any hostname supplied in the post body against the server blacklist.
         *
         * @param hostname a hostname to check against.
         * @return a blacklist result.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.POST)
        public BlacklistResult checkBody(@Nonnull @RequestBody String hostname) {
                return this.check(hostname);
        }

        /**
         * Checks a hostname against the blacklist.
         *
         * @param hostname a hostname.
         * @return a blacklist result.
         */
        @Nonnull
        private BlacklistResult checkHostname(@Nonnull String hostname) {
                List<String> hostnameParts = Splitter.on('.').splitToList(hostname);

                for (int i = 1; i < hostnameParts.size(); ++i) {
                        String currentHostname = "*." + Joiner.on('.').join(hostnameParts.subList(i, hostnameParts.size()));
                        String hash = Hashing.sha1().hashString(currentHostname, StandardCharsets.ISO_8859_1).toString();

                        if (this.hashes.contains(hash)) {
                                return new BlacklistResult(currentHostname, true);
                        }
                }

                return new BlacklistResult(hostname, false);
        }

        /**
         * Updates the local blacklist cache.
         */
        @PostConstruct
        @Scheduled(cron = "0 0 * * * ?")
        public void updateHashes() {
                try {
                        logger.info("Updating blacklist cache ...");
                        URL url = new URL(BLACKLIST_URL);

                        try (InputStream inputStream = url.openStream()) {
                                try (Reader reader = new InputStreamReader(inputStream)) {
                                        try (BufferedReader bufferedReader = new BufferedReader(reader)) {
                                                this.hashes.clear();
                                                bufferedReader.lines().forEach(this.hashes::add);
                                                logger.info("Found %d blacklist entries.", this.hashes.size());
                                        }
                                }
                        }
                } catch (IOException ex) {
                        logger.error("Could not update blacklist cache: %s", ex.getMessage());
                }
        }
}
