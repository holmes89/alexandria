(ns links.middleware
  (:require
   [environ.core :refer [env]]
   [buddy.auth :refer [authenticated?]]
   [buddy.auth.middleware :refer [wrap-authentication]]
   [buddy.auth.backends :refer [jws]]
   [ring.util.http-response :refer [unauthorized]]
   [clojure.tools.logging :as log]))

(def token-backend
  (jws {:secret (env :secret) :options {:alg :hs256}}))

(defn authenticated
  [handler]
  (fn [request]
    (if (authenticated? request)
      (handler request)
      (unauthorized {:error "Not authorized"}))))

(defn token-auth
  [handler]
  (wrap-authentication handler token-backend))

(defn logging
  [handler]
  (fn [request]
    (log/info request)
    (handler request)))

(defn cors
  "Cross-origin Resource Sharing (CORS) middleware. Allow requests from all
   origins, all http methods and Authorization and Content-Type headers."
  [handler]
  (fn [request]
    (let [response (handler request)]
      (-> response
          (assoc-in [:headers "Access-Control-Allow-Origin"] "*")
          (assoc-in [:headers "Access-Control-Allow-Methods"] "GET, PUT, PATCH, POST, DELETE, OPTIONS")
          (assoc-in [:headers "Access-Control-Allow-Headers"] "Authorization, Content-Type")))))
