Feature: Reconcile AddressableServices

    Scenario Outline: Reconciling <key> causes <result>.

        Given the following objects:
            """
            """
        And an AddressableService reconciler
        When reconciling "<key>"
        Then expect <result>

        Examples:
            | key            | result  |
            | too/many/parts | nothing |
            | foo/not-found  | nothing |

    # -----------------------------------------

    Scenario Outline: Reconciling but missing the service for generation <generation>.


        Given the following objects:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: <generation>
            spec:
              serviceName: webhook
            status:
              observedGeneration: 0
            """
        And an AddressableService reconciler

        When reconciling "ns/rut"

        Then expect status updates:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: <generation>
            spec:
              serviceName: webhook
            status:
              observedGeneration: <generation>
              conditions:
              - type: Ready
                status: "False"
                reason: ServiceUnavailable
                message: Service "webhook" wasn't found.
            """
        And expect Kubernetes Events:
            | Type   | Reason                       | Message                                 |
            | Normal | AddressableServiceReconciled | AddressableService reconciled: "ns/rut" |

        Examples:
            | generation |
            | 0          |
            | 1          |
            | 2          |

    # -----------------------------------------

    Scenario: Update status.address on spec.serviceName update.

        Given the following objects:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: 2
            spec:
              serviceName: webhook
            status:
              observedGeneration: 1
              address:
                url: http://old-webhook.ns.svc.cluster.local
              conditions:
              - type: Ready
                status: "True"
            ---
            apiVersion: v1
            kind: Service
            metadata:
              name: webhook
              namespace: ns
            spec:
              clusterIP: 10.20.30.40
              ports:
              - name: http
                port: 80
                protocol: TCP
                targetPort: 8080
              sessionAffinity: None
              type: ClusterIP
            """
        And an AddressableService reconciler

        When reconciling "ns/rut"

        Then expect status updates:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: 2
            spec:
              serviceName: webhook
            status:
              observedGeneration: 2
              address:
                url: http://webhook.ns.svc.cluster.local
              conditions:
              - type: Ready
                status: "True"
            """
        And expect Kubernetes Events:
            | Type   | Reason                       | Message                                 |
            | Normal | AddressableServiceReconciled | AddressableService reconciled: "ns/rut" |

    # -----------------------------------------

    Scenario Outline: Reconciling Normally for generation <generation>.

        Given the following objects:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: <generation>
            spec:
              serviceName: webhook
            ---
            apiVersion: v1
            kind: Service
            metadata:
              name: webhook
              namespace: ns
            spec:
              clusterIP: 10.20.30.40
              ports:
              - name: http
                port: 80
                protocol: TCP
                targetPort: 8080
              sessionAffinity: None
              type: ClusterIP
            """
        And an AddressableService reconciler

        When reconciling "ns/rut"

        Then expect status updates:
            """
            apiVersion: samples.knative.dev/v1alpha1
            kind: AddressableService
            metadata:
              name: rut
              namespace: ns
              generation: <generation>
            spec:
              serviceName: webhook
            status:
              observedGeneration: <generation>
              address:
                url: http://webhook.ns.svc.cluster.local
              conditions:
              - type: Ready
                status: "True"
            """
        And expect Kubernetes Events:
            | Type   | Reason                       | Message                                 |
            | Normal | AddressableServiceReconciled | AddressableService reconciled: "ns/rut" |

        Examples:
            | generation |
            | 0          |
            | 1          |
            | 2          |
