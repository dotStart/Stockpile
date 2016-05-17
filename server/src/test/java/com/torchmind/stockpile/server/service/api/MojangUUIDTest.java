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
package com.torchmind.stockpile.server.service.api;

import org.junit.Assert;
import org.junit.Test;

import java.util.UUID;

/**
 * <strong>Mojang UUID Test</strong>
 *
 * Tests the Mojang UUID representation against the relevant specifications.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 * @see MojangUUID
 */
public class MojangUUIDTest {
        private static final String MOJANG_UUID_STRING = "4566e69fc90748ee8d71d7ba5aa00d20";
        private static final MojangUUID MOJANG_UUID_POJO = new MojangUUID(MOJANG_UUID_STRING);
        private static final String UUID_STRING = "4566e69f-c907-48ee-8d71-d7ba5aa00d20";
        private static final UUID UUID_POJO = UUID.fromString(UUID_STRING);

        /**
         * Tests the conversion from Mojang UUIDs to strings.
         *
         * @see MojangUUID#toString()
         */
        @Test
        public void testToString() {
                Assert.assertEquals(MOJANG_UUID_STRING, MOJANG_UUID_POJO.toString());
        }

        /**
         * Tests the conversion from Mojang UUIDs to regular UUIDs.
         *
         * @see MojangUUID#toUUID()
         */
        @Test
        public void testToUUID() {
                Assert.assertEquals(UUID_POJO, MOJANG_UUID_POJO.toUUID());
        }

}
