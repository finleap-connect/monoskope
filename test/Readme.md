# Extra Test Formats

The tests in this folder are place here in order to separate test mechanisms and ensure a fast unit 
test is available for developers.

The unit tests are colocated with the sources they test and the integration test can be found in the `internal/test` folder. These can be tested using the `go test` command.

Test in this folder either require extra or external set up or use different tools or additional tools to work and may take considerable time to complete. They are not intented to be called constantly during development.

## Acceptance Tests

These tests use requirements stated in Gherkin to ensure the business requirements are continued to be fullfilled by the test. They use the `godog` utility to be evaluated. See the (godog github repository)[https://github.com/cucumber/godog] for further information.
