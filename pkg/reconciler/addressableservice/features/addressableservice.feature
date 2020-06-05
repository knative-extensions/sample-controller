Feature: Reconcile AddressableServices

    Scenario Outline: Reconciling <key> causes <result>.

    Given the following objects:
        """
        apiVersion: samples.knative.dev/v1alpha1
        kind: AddressableService
        metadata:
          name: make-it-addressable
          namespace: knative-samples
        spec:
          serviceName: webhook
        """

        And an AddressableService reconciler

    When reconciling "<key>"

    Then expect <result>

        And an unmodified cache

    Examples:
        | key               | result  |
        | too/many/parts  | nothing |
        | foo/not-found   | nothing |
