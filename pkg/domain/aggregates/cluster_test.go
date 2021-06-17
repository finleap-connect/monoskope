package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	meta "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedName                = "the one cluster"
	expectedLabel               = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
)

var _ = Describe("Unit Test for Cluster Aggregate", func() {

	var (
		expectedJWT = "thisisnotajwt"
		expectedCSR = []byte("This should be a CSR")
	)

	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New(), NewTestAggregateManager())

		err = createCluster(ctx, agg)
		Expect(err).NotTo(HaveOccurred())

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.ClusterCreated))

		data := &eventdata.ClusterCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedName))
		Expect(data.Label).To(Equal(expectedLabel))
		Expect(data.ApiServerAddress).To(Equal(expectedApiServerAddress))
		Expect(data.CaCertificateBundle).To(Equal(expectedClusterCACertBundle))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New(), NewTestAggregateManager())
		agg.IncrementVersion()

		ed := es.ToEventDataFromProto(&eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		})
		esEvent := es.NewEvent(ctx, events.ClusterCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).name).To(Equal(expectedName))
		Expect(agg.(*ClusterAggregate).label).To(Equal(expectedLabel))
		Expect(agg.(*ClusterAggregate).apiServerAddr).To(Equal(expectedApiServerAddress))
		Expect(agg.(*ClusterAggregate).caCertBundle).To(Equal(expectedClusterCACertBundle))

	})

	It("should write the jwt from a BootstrapTokenCreated event to the aggregate", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New(), NewTestAggregateManager())
		agg.IncrementVersion()

		ed := es.ToEventDataFromProto(&eventdata.ClusterBootstrapTokenCreated{
			JWT: expectedJWT,
		})
		esEvent := es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).bootstrapToken).To(Equal(expectedJWT))
	})

	It("should handle a ClusterCertificateRequested event", func() {
		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New(), NewTestAggregateManager())
		agg.IncrementVersion()

		ed := es.ToEventDataFromProto(&eventdata.ClusterCertificateRequested{
			CertificateSigningRequest: expectedCSR,
		})
		esEvent := es.NewEvent(ctx, events.ClusterCertificateRequested, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).certificateSigningRequest).To(Equal(expectedCSR))
	})

	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New(), NewTestAggregateManager())

		err = createCluster(ctx, agg)
		Expect(err).NotTo(HaveOccurred())

		esCommand, ok := cmd.NewRequestClusterCertificateCommand(uuid.New()).(*cmd.RequestClusterCertificateCommand)
		Expect(ok).To(BeTrue())

		esCommand.CertificateSigningRequest = expectedCSR

		err = agg.HandleCommand(ctx, esCommand)
		Expect(err).NotTo(HaveOccurred())

		event := agg.UncommittedEvents()[1]

		Expect(event.EventType()).To(Equal(events.ClusterCertificateRequested))

		data := &eventdata.ClusterCertificateRequested{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.CertificateSigningRequest).To(Equal(expectedCSR))

	})
})

func createCluster(ctx context.Context, agg es.Aggregate) error {
	esCommand, ok := cmd.NewCreateClusterCommand(uuid.New()).(*cmd.CreateClusterCommand)
	Expect(ok).To(BeTrue())

	esCommand.CreateCluster.Name = expectedName
	esCommand.CreateCluster.Label = expectedLabel
	esCommand.CreateCluster.ApiServerAddress = expectedApiServerAddress
	esCommand.CreateCluster.ClusterCACertBundle = expectedClusterCACertBundle

	return agg.HandleCommand(ctx, esCommand)
}

func makeMetadataContextWithSystemAdminUser() (context.Context, error) {
	metaMgr, err := meta.NewDomainMetadataManager(context.Background())
	if err != nil {
		return nil, err
	}

	// forces the setting of the domain context
	metaMgr, err = meta.NewDomainMetadataManager(metaMgr.GetContext())
	if err != nil {
		return nil, err
	}

	metaMgr.SetUserInformation(&meta.UserInformation{
		Id:     uuid.New(),
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})

	metaMgr.SetRoleBindings([]*projections.UserRoleBinding{
		{
			Role:  roles.Admin.String(),
			Scope: scopes.System.String(),
		},
	})

	return metaMgr.GetContext(), nil

}

type aggregateTestStore struct {
}

// NewTestAggregateManager creates a new dummy AggregateHandler which allows observing interactions and injecting test data.
func NewTestAggregateManager() es.AggregateStore {
	return &aggregateTestStore{}
}

// Get returns the most recent version of all aggregate of a given type.
func (tas *aggregateTestStore) All(context.Context, es.AggregateType) ([]es.Aggregate, error) {
	return []es.Aggregate{}, nil
}

// Get returns the most recent version of an aggregate.
func (tas *aggregateTestStore) Get(context.Context, es.AggregateType, uuid.UUID) (es.Aggregate, error) {
	return nil, nil
}

// Update stores all in-flight events for an aggregate.
func (tas *aggregateTestStore) Update(context.Context, es.Aggregate) error {
	return nil
}
