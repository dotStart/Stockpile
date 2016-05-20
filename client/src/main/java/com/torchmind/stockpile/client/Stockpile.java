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
package com.torchmind.stockpile.client;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.torchmind.stockpile.data.v1.BlacklistResult;
import com.torchmind.stockpile.data.v1.NameLookupResult;
import com.torchmind.stockpile.data.v1.PlayerProfile;
import com.torchmind.stockpile.data.v1.ServerInformation;
import retrofit.Call;
import retrofit.JacksonConverterFactory;
import retrofit.Retrofit;
import retrofit.http.*;

import javax.annotation.Nonnull;
import java.util.UUID;

/**
 * <strong>Stockpile Client</strong>
 *
 * Provides a simple method of accessing the Stockpile API.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public interface Stockpile {

        /**
         * Creates a client instance for the specified base URL.
         *
         * @param baseUrl a base URL.
         * @return a client instance.
         */
        @Nonnull
        static Stockpile create(@Nonnull String baseUrl) {
                ObjectMapper mapper = new ObjectMapper();
                mapper.findAndRegisterModules();

                Retrofit retrofit = (new Retrofit.Builder())
                        .addConverterFactory(JacksonConverterFactory.create(mapper))
                        .baseUrl(baseUrl)
                        .build();

                return retrofit.create(Stockpile.class);
        }

        /**
         * Retrieves the current server and API version.
         *
         * @return a version.
         */
        @Nonnull
        @GET("/v1")
        Call<ServerInformation> getServerInformation();

        /**
         * Instructs the server to handle a server handshake on behalf of the requesting server in order to update local
         * profile records.
         *
         * @param username a username.
         * @param serverId a server identifier.
         * @return a player profile.
         */
        @Nonnull
        @GET("/v1/login")
        Call<PlayerProfile> login(@Nonnull @Query(value = "username", encoded = true) String username, @Nonnull @Query(value = "serverId", encoded = true) String serverId);

        /**
         * Checks a hostname against the server blacklist.
         *
         * @param hostname a hostname.
         * @return a blacklist result.
         */
        @Nonnull
        @POST("/v1/blacklist")
        Call<BlacklistResult> lookupBlacklistEntry(@Nonnull @Query(value = "hostname", encoded = true) String hostname);

        /**
         * Looks up the UUID which corresponds to the supplied name.
         *
         * @param name a name.
         * @return a lookup result.
         */
        @Nonnull
        @GET("/v1/name/{name}")
        Call<NameLookupResult> lookupName(@Nonnull @Path(value = "name", encoded = true) String name);

        /**
         * Looks up a profile based on its identifier.
         *
         * @param identifier an identifier.
         * @return a profile.
         */
        @Nonnull
        @GET("/v1/profile/{identifier}")
        Call<PlayerProfile> lookupProfile(@Nonnull @Path(value = "identifier", encoded = true) UUID identifier);

        /**
         * Looks up a profile based on its name.
         *
         * @param name a name.
         * @return a profile.
         */
        @Nonnull
        @GET("/v1/profile/{name}")
        Call<PlayerProfile> lookupProfile(@Nonnull @Path(value = "name", encoded = true) String name);

        /**
         * Purges a display name from the server.
         *
         * @param name a name.
         */
        @DELETE("/v1/name/{name}")
        void purgeName(@Nonnull @Path(value = "name", encoded = true) String name);

        /**
         * Purges a player profile from the server.
         *
         * @param identifier a profile identifier.
         */
        @DELETE("/v1/profile/{identifier}")
        void purgeProfile(@Nonnull @Path(value = "identifier", encoded = true) UUID identifier);

        /**
         * Purges a player profile from the server using its display name.
         *
         * @param name a name.
         */
        @DELETE("/v1/profile/{name}")
        void purgeProfile(@Nonnull @Path(value = "name", encoded = true) String name);

        /**
         * Requests the server to shut down gracefully.
         */
        @POST("/v1/shutdown")
        void requestShutdown();
}
