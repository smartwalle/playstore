package playstore

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

type options struct {
	client *http.Client
}

type Option func(opts *options)

func WithClient(client *http.Client) Option {
	return func(opts *options) {
		if client != nil {
			opts.client = client
		}
	}
}

type Client struct {
	products        *androidpublisher.PurchasesProductsService
	subscriptions   *androidpublisher.PurchasesSubscriptionsService
	subscriptionsV2 *androidpublisher.PurchasesSubscriptionsv2Service
	voidedPurchases *androidpublisher.PurchasesVoidedpurchasesService
}

func New(jsonKey []byte, opts ...Option) (*Client, error) {
	var nOpt = &options{
		client: http.DefaultClient,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(nOpt)
		}
	}

	var ctx = context.WithValue(context.Background(), oauth2.HTTPClient, nOpt.client)
	var conf, err = google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	var val = conf.Client(ctx).Transport.(*oauth2.Transport)
	if _, err = val.Source.Token(); err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	var products = androidpublisher.NewPurchasesProductsService(service)
	var subscriptions = androidpublisher.NewPurchasesSubscriptionsService(service)
	var subscriptionsV2 = androidpublisher.NewPurchasesSubscriptionsv2Service(service)
	var voidedPurchases = androidpublisher.NewPurchasesVoidedpurchasesService(service)

	return &Client{
		products:        products,
		subscriptions:   subscriptions,
		subscriptionsV2: subscriptionsV2,
		voidedPurchases: voidedPurchases,
	}, nil
}

func NewWithJSONFile(filename string, opts ...Option) (*Client, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return New(data, opts...)
}

func (this *Client) ProductsService() *androidpublisher.PurchasesProductsService {
	return this.products
}

func (this *Client) SubscriptionsService() *androidpublisher.PurchasesSubscriptionsService {
	return this.subscriptions
}

func (this *Client) SubscriptionsV2Service() *androidpublisher.PurchasesSubscriptionsv2Service {
	return this.subscriptionsV2
}

func (this *Client) VoidedPurchasesService() *androidpublisher.PurchasesVoidedpurchasesService {
	return this.voidedPurchases
}

func (this *Client) VerifyProduct(ctx context.Context, packageName, productId, token string) (*androidpublisher.ProductPurchase, error) {
	return this.products.Get(packageName, productId, token).Context(ctx).Do()
}

func (this *Client) AcknowledgeProduct(ctx context.Context, packageName, productID, token, developerPayload string) error {
	var param = &androidpublisher.ProductPurchasesAcknowledgeRequest{DeveloperPayload: developerPayload}
	return this.products.Acknowledge(packageName, productID, token, param).Context(ctx).Do()
}

func (this *Client) ConsumeProduct(ctx context.Context, packageName, productId, token string) error {
	return this.products.Consume(packageName, productId, token).Context(ctx).Do()
}

func (this *Client) VerifySubscription(ctx context.Context, packageName, subscriptionId, token string) (*androidpublisher.SubscriptionPurchase, error) {
	return this.subscriptions.Get(packageName, subscriptionId, token).Context(ctx).Do()
}

func (this *Client) VerifySubscriptionV2(ctx context.Context, packageName, token string) (*androidpublisher.SubscriptionPurchaseV2, error) {
	return this.subscriptionsV2.Get(packageName, token).Context(ctx).Do()
}

func (this *Client) AcknowledgeSubscription(ctx context.Context, packageName, subscriptionId, token, developerPayload string) error {
	var param = &androidpublisher.SubscriptionPurchasesAcknowledgeRequest{DeveloperPayload: developerPayload}
	return this.subscriptions.Acknowledge(packageName, subscriptionId, token, param).Context(ctx).Do()
}

func (this *Client) CancelSubscription(ctx context.Context, packageName string, subscriptionId string, token string) error {
	return this.subscriptions.Cancel(packageName, subscriptionId, token).Context(ctx).Do()
}

func (this *Client) RefundSubscription(ctx context.Context, packageName string, subscriptionId string, token string) error {
	return this.subscriptions.Refund(packageName, subscriptionId, token).Context(ctx).Do()
}

func (this *Client) RevokeSubscription(ctx context.Context, packageName string, subscriptionId string, token string) error {
	return this.subscriptions.Revoke(packageName, subscriptionId, token).Context(ctx).Do()
}
