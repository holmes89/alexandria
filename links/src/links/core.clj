(ns links.core
  (:require [ring.adapter.jetty :refer [run-jetty]]
            [toucan.db :as db]
            [toucan.models :as models]
            [compojure.api.sweet :refer [api]]
            [links.link :refer [link-routes]]
            [environ.core :refer [env]])
  (:gen-class))

(def db-spec
  {:classname   "org.postgresql.Driver"
   :connection-uri (env :database-url "//db:5432/mind")})

(def app (api link-routes))

(defn -main
  [& args]
  (db/set-default-db-connection! db-spec)
  (models/set-root-namespace! 'links.models)
  (run-jetty app {:port 3000}))
