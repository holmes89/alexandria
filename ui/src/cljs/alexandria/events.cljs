(ns alexandria.events
  (:require
   [ajax.core :as ajax]        
   [day8.re-frame.http-fx]
   [re-frame.core :as re-frame]
   [alexandria.db :as db]
   ))

(re-frame/reg-event-db
 ::initialize-db
 (fn [_ _]
   db/default-db))

(re-frame/reg-event-db
    ::set-active-panel
  (fn [db [_ active-panel]]
    (assoc db :active-panel active-panel)))

(re-frame/reg-event-db                   
    ::process-response             
  (fn
    [db [_ response]]
    (-> db
        (assoc :loading? false)
        (assoc :document-data (js->clj response)))))

(re-frame/reg-event-db                   
    ::bad-response             
  (fn
    [message]
    (js/console.log  message)))

(re-frame/reg-event-fx
    ::get-documents
  (fn
    [{db :db} _]
    {:http-xhrio {:method          :get
                  :uri             "http://localhost:8080/documents/"
                  :format          (ajax/json-request-format)
                  :response-format (ajax/json-response-format {:keywords? true}) 
                  :on-success      [::process-response]
                  :on-failure      [::bad-response]
                  }
     :db  (assoc db :loading? true)}))
