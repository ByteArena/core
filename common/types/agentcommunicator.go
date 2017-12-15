package types

type AgentCommunicatorInterface interface {
	NetSenderInterface
	AgentMutationBatcherInterface
}
