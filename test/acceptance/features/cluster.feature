Feature: Clusters
   Cluster aggregates hold information about Kubernetes clusters that are accessable and managable by the
   Monoskope control plane

   Scenario: create a cluster
     Given there is an empty list of clusters
     And my name is "admin", my email is "admin@monoskope.io" and have a token issued by "monoskope"
     
     When I create a cluster with the dns address "one.example.com"
     
     Then there should be a cluster with the name "one_example_com" in the list of clusters
     # the user is needed so the M8 operator can authenticate itself when registering
     # see https://finleap-connect.atlassian.net/browse/FCLOUD-4050?focusedCommentId=219638
     And there should be a user with the name "one.example.com" in the list of users
     And there should be a role binding for the user "one.example.com" and the role "user" for the scope "cluster" and resource "one.example.com"


   Scenario: cluster dns addresses must be unique
     Given there are clusters with the dns addresses of:
     And my name is "admin", my email is "admin@monoskope.io" and have a token issued by "monoskope"
     | one.example.com |
     When I create a cluster with the dns address "one.example.com"
     Then the command should fail.
  
