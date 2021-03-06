package weatherservice

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	weatherservicev1alpha1 "github.com/embano1/weatherservice/pkg/apis/weatherservice/v1alpha1"
	"github.com/embano1/weatherservice/pkg/internal/openweather"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var appIDFile string

func init() {
	flag.StringVar(&appIDFile, "c", "", "Credentials file for the OpenWeather API APPID")
}

// Add creates a new WeatherService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {

	if appIDFile == "" {
		log.Fatalln("Please specify a credentials file for the OpenWeather API APPID")
	}

	appID, err := ioutil.ReadFile(appIDFile)
	if err != nil {
		log.Fatalf("could not read credentials file: %v", err)
	}

	owc := openweather.NewClient(string(appID), false)

	return &ReconcileWeatherService{
		client:            mgr.GetClient(),
		scheme:            mgr.GetScheme(),
		OpenWeatherClient: owc,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("weatherservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource WeatherService
	err = c.Watch(&source.Kind{Type: &weatherservicev1alpha1.WeatherService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner WeatherService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &weatherservicev1alpha1.WeatherService{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileWeatherService{}

// ReconcileWeatherService reconciles a WeatherService object
type ReconcileWeatherService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client            client.Client
	scheme            *runtime.Scheme
	OpenWeatherClient *openweather.Client
}

// Reconcile reads that state of the cluster for a WeatherService object and makes changes based on the state read
// and what is in the WeatherService.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWeatherService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling WeatherService %s/%s\n", request.Namespace, request.Name)

	// Fetch the WeatherService instance
	instance := &weatherservicev1alpha1.WeatherService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	city := instance.Spec.City
	metric := instance.Spec.Unit

	res, err := r.OpenWeatherClient.Get(city, metric)
	if err != nil {
		// TODO: print but dont tell it's an error for now
		log.Println(err)
		return reconcile.Result{}, nil
	}

	instance.Status.Temperature = res.Main.Temp
	instance.Status.Unit = metric
	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
