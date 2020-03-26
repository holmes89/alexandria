(ns links.link
  (:require [schema.core :as s]
            [links.models.link :refer [Link]]
            [toucan.db :as db]
            [ring.util.http-response :refer [ok not-found created]]
            [compojure.api.sweet :refer :all]
            [links.middleware :as mw]
            [clojure.tools.logging :as log]))

(s/defschema LinkSchema
  {(s/optional-key :id) s/Uuid
   :link s/Str
   :display_name s/Str
   (s/optional-key :description) s/Str
   (s/optional-key :created) s/Str
   })

(defn response [link]
  (if link
    (ok link)
    (not-found)))

(defn create-link-handler [create-link-req]
  (-> (db/insert! Link create-link-req)
      (response)))

(defn update-link-handler [id update-link-req]
  (db/update! Link id update-link-req)
  (ok))

(defn get-link-handler [id]
  (-> (Link id)
      (response)))

(defn find-all-handler []
  (-> (db/select Link)
      (response)))

(defroutes link-routes
  (context "/links" []
           :middleware [mw/cors]
           (POST "/" request
             :responses {created LinkSchema}
             :body [create-link-req LinkSchema]
             (create-link-handler create-link-req))
           (GET "/" []
             (find-all-handler))
           (GET "/:id" []
             :path-params [id :- s/Uuid]
             (get-link-handler id))
           (PATCH "/:id" request
             :path-params [id :- s/Uuid]
             :body [update-link-req LinkSchema]
             (update-link-handler id update-link-req request))))
