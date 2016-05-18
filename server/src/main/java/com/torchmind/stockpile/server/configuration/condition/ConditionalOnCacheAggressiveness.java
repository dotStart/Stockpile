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
import org.springframework.context.annotation.Conditional;

import java.lang.annotation.*;

/**
 * <strong>Conditional on Cache Aggressiveness</strong>
 *
 * Instructs Spring to only initialize a bean if the configured application aggressiveness is equal to a certain preset.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Documented
@Retention(RetentionPolicy.RUNTIME)
@Conditional(CacheAggressivenessCondition.class)
@Target({ElementType.TYPE, ElementType.METHOD})
public @interface ConditionalOnCacheAggressiveness {
        Aggressiveness value() default Aggressiveness.UNKNOWN;
}
