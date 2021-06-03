Feature: Clusters
  Cluster aggregates hold information about Kubernetes clusters that are accessable and managable by the
  Monoskope control plane.
  Clusters are also used to represent the local state kept for keeping track of registered M8 Operators
  running in remote clusters.
  In general, cluster information should only be manipulated by admins with the system scope.


  Scenario: attempting to create a cluster with insufficient credentials
    Given there is an empty list of clusters
    And my name is "not-an-admin", my email is "not-an-admin@monoskope.io" and have a token issued by "monoskope"
     
    When I create a cluster with the dns address "one.example.com"

    Then the command should fail.


  Scenario: attempting to create a cluster with an exaisting name
    Given there are clusters with the names of:
     | any-one-cluster |
    And my name is "admin", my email is "admin@monoskope.io" and have a token issued by "monoskope"
     
    When I create a cluster with the dns address "one.example.com" and the name "any-one-cluster"

    Then the command should fail.


  Scenario: cluster dns addresses must be unique
    Given there are clusters with the dns addresses of:
     | one.example.com |
    And my name is "admin", my email is "admin@monoskope.io" and have a token issued by "monoskope"
     
    When I create a cluster with the dns address "one.example.com"

    Then the command should fail.


  Scenario: create a cluster
    Given there is an empty list of clusters
    And my name is "admin", my email is "admin@monoskope.io" and have a token issued by "monoskope"
     
    When I create a cluster with the dns address "one.example.com"
     
    Then there should be a cluster with the dns address "one.example.com" in the list of clusters
    # the user is needed so the M8 operator can authenticate itself when registering
    # see https://finleap-connect.atlassian.net/browse/FCLOUD-4050?focusedCommentId=219638
    And there should be a user with the name "one.example.com" in the list of users
    And there should be a role binding for the user "one.example.com" and the role "user" for the scope "cluster" and resource "one.example.com"
    And there should be a JWT token available for the cluster that is valid for the the name "one.example.com"

  
