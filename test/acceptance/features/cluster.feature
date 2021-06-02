Feature: Clusters
   Cluster aggregates hold information about Kubernetes clusters that are accessable and managable by the
   Monoskope control plane

   Scenario: create a cluster
     Given there is an empty list of clusters
     When I create a cluster with the dns address "one.example.com"
     Then there should be a cluster with the name "one_example_com" in the list of clusters.


   Scenario: cluster dns addresses must be unique
     Given there are clusters with the dns addresses of:
     | one.example.com |
     When I create a cluster with the dns address "one.example.com"
     Then the command should fail.
  