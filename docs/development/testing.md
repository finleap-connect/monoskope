# Testing

There are three general areas of testing within the code:

* **Unit Tests**, that are co-located with the functions that they are testing. These are implemented using [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) to aid in readability. These should be implemented using TDD and BDD principles.
* **Integration Tests**, the reconstruct the complete software stack automatically. These should be used as the primary test environment for developers to verify that new modules fit with the rest of the system. They are also implemented using Ginkgo and Gomega.
* **Acceptance Tests**, they ensure that the business rules are correctly implemented. They are written in Gherkin and use [godog](https://github.com/cucumber/godog) to validate agains the code.


