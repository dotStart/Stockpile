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

import com.torchmind.stockpile.data.v1.ServerInformation;
import com.torchmind.stockpile.data.v1.Version;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Nonnull;
import javax.annotation.concurrent.ThreadSafe;

/**
 * <strong>Root Controller</strong>
 *
 * Handles requests to the root endpoint within the API which informs clients of the local application version.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@ThreadSafe
@RestController
@RequestMapping("/v1")
public class ApplicationController {
        public static final Version VERSION = new Version(1, Version.State.DEVELOPMENT);
        private final ApplicationContext applicationContext;

        @Autowired
        public ApplicationController(@Nonnull ApplicationContext applicationContext) {
                this.applicationContext = applicationContext;
        }

        /**
         * <code>/v1/</code>
         *
         * Returns version information on this specific API endpoint.
         *
         * @return a version information representation.
         */
        @Nonnull
        @RequestMapping(method = RequestMethod.GET)
        public ServerInformation get() {
                Package p = this.getClass().getPackage();

                if (p != null) {
                        return new ServerInformation(p, VERSION);
                }

                return new ServerInformation(VERSION);
        }

        /**
         * <code>/v1/shutdown</code>
         *
         * Shuts down the server.
         */
        @RequestMapping(path = "/shutdown", method = RequestMethod.POST)
        public void shutdown() {
                if (!(this.applicationContext instanceof ConfigurableApplicationContext)) {
                        throw new IllegalStateException("Cannot shut down context: Unknown implementation");
                }

                (new Thread() {
                        @Override
                        public void run() {
                                // wait a bit to ensure the request has finished before we shut down
                                try {
                                        Thread.sleep(500L);
                                } catch (InterruptedException ignore) {
                                }

                                // actually instruct the context to shut down
                                ((ConfigurableApplicationContext) applicationContext).close();
                        }
                }).start();
        }
}
