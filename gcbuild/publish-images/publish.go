/* Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

// We'd have an array here for packages without multiple images.

var sharedPackages = map[string][]string {
	"pubsub-binary-to-bigquery": {"pubsub-proto-to-bigquery", "pubsub-avro-to-bigquery"},
}

func publish(packageName string, image string) error {
	args := []string {
		"clean", "package", "-f", "v2/pom.xml",
		"-Dmaven.test.skip",
		// In reality, the project would be a flag.
		fmt.Sprintf("-Dimage=gcr.io/zhoufek-test-331019/rc/%s", image),
		"-Dbase-container-image=gcr.io/dataflow-templates-base/java8-template-launcher-base",
		"-Dbase-container-image.version=latest",
		fmt.Sprintf("-Dapp-root=/template/%s", image),
		fmt.Sprintf("-Dcommand-spec=/template/%[1]s/resources/%[1]s-command-spec.json", image),
		"-pl", packageName, "-am", "-ntp",
	}

	cmd := exec.Command("mvn", args...)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return cmd.Wait()
}

func main() {
	// Realistically, the first step will be to check if the images already exist
	// in GCR. Packaging and pushing them is relatively expensive (multiple minutes).

	for k, v := range sharedPackages {
		for _, image := range v {
			if err := publish(k, image); err != nil {
				log.Fatalf("Failed when publishing %v: %v", image, err)
			}
		}
	}
}