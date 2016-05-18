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

import com.torchmind.stockpile.data.v1.NameLookupResult;
import com.torchmind.stockpile.server.service.api.ProfileService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Nonnull;

/**
 * <strong>Name Controller</strong>
 *
 * Provides methods for looking up display names.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@RestController
@RequestMapping("/v1/name")
public class NameController {
        private final ProfileService profileService;

        @Autowired
        public NameController(@Nonnull ProfileService profileService) {
                this.profileService = profileService;
        }

        /**
         * <code>GET /v1/name/{name}</code>
         *
         * Attempts to find the corresponding profile identifier for the specified name.
         *
         * @param name a name.
         * @return a lookup result.
         */
        @Nonnull
        @RequestMapping(path = "/{name}", method = RequestMethod.GET)
        public NameLookupResult lookup(@Nonnull @PathVariable("name") String name) {
                return new NameLookupResult(this.profileService.findIdentifier(name));
        }

        /**
         * <code>POST /v1/name</code>
         *
         * Attempts to find the corresponding profile identifier for the specified name.
         *
         * @param name a name.
         * @return a lookup result.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.POST)
        public NameLookupResult lookupBody(@Nonnull @RequestBody String name) {
                return this.lookup(name);
        }
}
