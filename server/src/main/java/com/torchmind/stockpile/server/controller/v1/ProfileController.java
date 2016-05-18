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

import com.torchmind.stockpile.data.v1.PlayerProfile;
import com.torchmind.stockpile.server.entity.Profile;
import com.torchmind.stockpile.server.service.api.ProfileService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Nonnull;
import javax.annotation.concurrent.ThreadSafe;
import java.util.UUID;
import java.util.regex.Pattern;

/**
 * <strong>Profile Controller</strong>
 *
 * Provides methods for accessing profile information.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@ThreadSafe
@RestController
@RequestMapping("/v1/profile")
public class ProfileController {
        public static final Pattern UUID_PATTERN = Pattern.compile("^[A-F0-9]{8}-[A-F0-9]{4}-[1-5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[A-F0-9]{12}$", Pattern.CASE_INSENSITIVE);
        private final ProfileService profileService;

        @Autowired
        public ProfileController(@Nonnull ProfileService profileService) {
                this.profileService = profileService;
        }

        /**
         * <code>GET /v1/profile/{name}</code>
         *
         * Looks up a profile based on its name or identifier.
         *
         * @param name a name or identifier.
         * @return a response entity.
         */
        @Nonnull
        @RequestMapping(path = "/{name}", method = RequestMethod.GET)
        public ResponseEntity<PlayerProfile> lookup(@Nonnull @PathVariable("name") String name) {
                final Profile profile;

                if (UUID_PATTERN.matcher(name).matches()) {
                        UUID identifier = UUID.fromString(name);
                        profile = this.profileService.get(identifier);
                } else {
                        profile = this.profileService.find(name);
                }

                return new ResponseEntity<>(profile.toRestRepresentation(), (profile.isCached() ? HttpStatus.OK : HttpStatus.CREATED));
        }

        /**
         * <code>POST /v1/profile</code>
         *
         * Looks up a profile based on its name or identifier.
         *
         * @param name a name or identifier.
         * @return a response entity.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.POST)
        public ResponseEntity<PlayerProfile> lookupBody(@Nonnull @RequestBody String name) {
                return this.lookup(name);
        }

        /**
         * <code>DELETE /v1/profile/{name}</code>
         *
         * Purges a profile from the cache.
         *
         * @param name a name or identifier.
         */
        @ResponseStatus(HttpStatus.NO_CONTENT)
        @RequestMapping(path = "/{name}", method = RequestMethod.DELETE)
        public void purge(@Nonnull @PathVariable("name") String name) {
                if (UUID_PATTERN.matcher(name).matches()) {
                        UUID identifier = UUID.fromString(name);
                        this.profileService.purge(identifier);
                        return;
                }

                this.profileService.purge(name);
        }

        /**
         * <code>DELETE /v1/profile</code>
         *
         * Purges a profile from the cache.
         *
         * @param name a name or identifier.
         */
        @ResponseStatus(HttpStatus.NO_CONTENT)
        @RequestMapping(method = RequestMethod.DELETE)
        public void purgeBody(@Nonnull @RequestBody String name) {
                this.purge(name);
        }
}
