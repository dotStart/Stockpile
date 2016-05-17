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
package com.torchmind.stockpile.server.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.web.AbstractErrorController;
import org.springframework.boot.autoconfigure.web.ErrorAttributes;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Nonnull;
import javax.servlet.http.HttpServletRequest;
import java.util.Map;

/**
 * <strong>Restful Error Controller</strong>
 *
 * Provides a controller which will handle all uncaught exceptions.
 *
 * @author <a href="mailto:johannesd@torchmind.com">Johannes Donath</a>
 */
@RestController
@RequestMapping("/error")
public class RestfulErrorController extends AbstractErrorController {

        @Autowired
        public RestfulErrorController(@Nonnull ErrorAttributes errorAttributes) {
                super(errorAttributes);
        }

        /**
         * Handles failed requests.
         *
         * @param request a request.
         * @return an error response.
         */
        @ResponseBody
        @RequestMapping
        public ResponseEntity<Map<String, Object>> error(HttpServletRequest request) {
                Map<String, Object> body = this.getErrorAttributes(request, false);
                HttpStatus status = this.getStatus(request);

                return new ResponseEntity<>(body, status);
        }

        /**
         * {@inheritDoc}
         */
        @Nonnull
        @Override
        public String getErrorPath() {
                return "/error";
        }
}
