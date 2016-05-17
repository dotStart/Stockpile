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
package com.torchmind.stockpile.server.entity.repository;

import com.torchmind.stockpile.server.entity.Profile;
import com.torchmind.stockpile.server.entity.ProfileProperty;
import org.springframework.data.repository.PagingAndSortingRepository;

import javax.annotation.Nonnull;
import java.util.Optional;
import java.util.UUID;

/**
 * <strong>Profile Property Repository</strong>
 *
 * Provides a management interface for finding, creating, updating and deleting profile properties.
 * Note: This interface is implemented by Spring during runtime.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public interface ProfilePropertyRepository extends PagingAndSortingRepository<ProfileProperty, UUID> {

        /**
         * Searches for a profile property with the specified name and profile.
         *
         * @param name    a property name.
         * @param profile a profile.
         * @return a profile property or, if no property within the specified profile was found, an empty optional.
         */
        @Nonnull
        Optional<ProfileProperty> findByNameAndProfile(@Nonnull String name, @Nonnull Profile profile);
}
