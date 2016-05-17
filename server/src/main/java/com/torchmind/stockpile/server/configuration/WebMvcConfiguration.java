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
package com.torchmind.stockpile.server.configuration;

import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.web.servlet.config.annotation.*;

import javax.annotation.Nonnull;

/**
 * <strong>Web MVC Configuration</strong>
 *
 * Configures Spring's web MVC component for use with pure REST based communication.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@Configuration
public class WebMvcConfiguration extends WebMvcConfigurerAdapter {

        /**
         * {@inheritDoc}
         */
        @Override
        public void addResourceHandlers(ResourceHandlerRegistry registry) {
                registry.addResourceHandler("/index.html").addResourceLocations("classpath:/static/index.html");
                registry.addResourceHandler("/logo.svg").addResourceLocations("classpath:/static/logo.svg");
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public void configureContentNegotiation(@Nonnull ContentNegotiationConfigurer configurer) {
                configurer.defaultContentType(MediaType.APPLICATION_JSON_UTF8);
                configurer.mediaType("json", MediaType.APPLICATION_JSON_UTF8);
                configurer.mediaType("xml", MediaType.APPLICATION_XML);
                configurer.favorPathExtension(true);
                configurer.ignoreUnknownPathExtensions(true);
                configurer.ignoreAcceptHeader(false);
        }

        /**
         * {@inheritDoc}
         */
        @Override
        public void addCorsMappings(@Nonnull CorsRegistry registry) {
                registry.addMapping("/**")
                        .allowedOrigins("*")
                        .allowedMethods("GET", "POST", "PUT", "DELETE")
                        .allowedHeaders("X-Authentication");
        }
}
