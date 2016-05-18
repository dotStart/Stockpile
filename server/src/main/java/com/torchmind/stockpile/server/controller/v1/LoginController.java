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
import com.torchmind.stockpile.server.service.api.MojangUUID;
import com.torchmind.stockpile.server.service.api.ProfileService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Nonnull;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * <strong>Login Controller</strong>
 *
 * Provides a proxy to Mojang's login endpoint to improve caching on compatible servers.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@RestController
@RequestMapping("/v1/login")
public class LoginController {
        private final ProfileService profileService;

        @Autowired
        public LoginController(@Nonnull ProfileService profileService) {
                this.profileService = profileService;
        }

        /**
         * <code>GET /v1/login</code>
         *
         * Proxies a login request and updates the cache.
         *
         * @param username a username.
         * @param serverId a server identifier.
         * @return a profile.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.GET)
        public PlayerProfile login(@Nonnull @RequestParam("username") String username, @Nonnull @RequestParam("serverId") String serverId) {
                return this.profileService.join(username, serverId).toRestRepresentation();
        }

        /**
         * <code>GET /v1/login</code>
         *
         * Proxies a login request and updates the cache.
         * Note: This version requires the X-Forward header to be present and causes the API to return a regular Mojang
         * response.
         *
         * @param username a username.
         * @param serverId a server identifier.
         * @return a profile.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.GET, produces = "application/json", headers = "X-Forward")
        public Map<String, Object> loginForward(@Nonnull @RequestParam("username") String username, @Nonnull @RequestParam("serverId") String serverId) {
                PlayerProfile profile = this.login(username, serverId);

                Map<String, Object> root = new HashMap<>();
                root.put("id", (new MojangUUID(profile.getIdentifier())).toString());
                root.put("name", profile.getName());

                {
                        List<Map<String, Object>> properties = new ArrayList<>();
                        profile.getProperties().forEach((p) -> {
                                Map<String, Object> property = new HashMap<>();
                                property.put("name", p.getName());
                                property.put("value", p.getValue());

                                if (p.getSignature() != null) {
                                        property.put("signature", p.getSignature());
                                }

                                properties.add(property);
                        });
                        root.put("properties", properties);
                }

                return root;
        }
}
