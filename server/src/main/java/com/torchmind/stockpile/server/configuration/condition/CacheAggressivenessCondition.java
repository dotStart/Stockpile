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
package com.torchmind.stockpile.server.configuration.condition;

import com.torchmind.stockpile.server.configuration.CacheConfiguration.Aggressiveness;
import org.springframework.context.annotation.Condition;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.core.type.AnnotatedTypeMetadata;

import javax.annotation.Nonnull;

/**
 * <strong>Cache Aggressiveness Condition</strong>
 *
 * Checks whether a specific cache aggressiveness is configured.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
public class CacheAggressivenessCondition implements Condition {

        /**
         * {@inheritDoc}
         */
        @Override
        public boolean matches(@Nonnull ConditionContext context, @Nonnull AnnotatedTypeMetadata metadata) {
                Aggressiveness aggressiveness = (Aggressiveness) metadata.getAnnotationAttributes(ConditionalOnCacheAggressiveness.class.getName()).getOrDefault("value", Aggressiveness.HIGH);
                return Aggressiveness.valueOf(context.getEnvironment().getProperty("cache.aggressiveness", "high").toUpperCase()) == aggressiveness;
        }
}
