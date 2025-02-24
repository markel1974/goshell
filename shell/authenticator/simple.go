/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authenticator

import (
	"crypto/sha512"
	"encoding/hex"
)

type SimpleAuthenticator struct {
	username string
	hash     string
	salt     []byte
}

func NewSimpleAuthenticator() *SimpleAuthenticator {
	a := &SimpleAuthenticator{
		username: "",
		hash:     "",
		salt:     []byte{},
	}
	return a
}

func (a *SimpleAuthenticator) Setup(username string) (string, error) {
	pwd, err := Generate(16)
	if err != nil {
		return "", err
	}
	salt, err := Generate(8)
	if err != nil {
		return "", err
	}
	a.salt = []byte(salt)
	a.hash = a.generateHash(pwd, a.salt)
	a.username = username
	return pwd, nil
}

func (a *SimpleAuthenticator) IsAuthenticated(user string, pass string) bool {
	if a.username != user {
		return false
	}
	hash := a.generateHash(pass, a.salt)
	return a.hash == hash
}

func (a *SimpleAuthenticator) generateHash(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}
