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
    ::next-page
  (fn [db]
    (update db :page-num inc 1)))

(re-frame/reg-event-db
    ::prev-page
  (fn [db]
    (update db :page-num dec 1)))

(re-frame/reg-event-db
    ::zoom-in
  (fn [db]
    (update db :zoom + 0.2)))

(re-frame/reg-event-db
    ::zoom-out
  (fn [db]
    (update db :zoom - 0.2)))

(re-frame/reg-event-db
    ::show-upload-modal
  (fn [db]
    (assoc db :show-upload true)))

(re-frame/reg-event-db
    ::hide-upload-modal
  (fn [db]
    (assoc db :show-upload false)))

(re-frame/reg-event-db
    ::update-upload-file-name
  (fn [db [_ name]]
    (assoc db :upload-file-name name)))

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

(re-frame/reg-event-db                   
    ::process-get-doc-by-id             
  (fn
    [db [_ response]]
    (-> db
        (assoc :loading? false)
        (assoc :active-doc (js->clj response)))))

(re-frame/reg-event-fx
    ::get-document-by-id
  (fn
    [{db :db} [_ id]]
    {:http-xhrio {:method          :get
                  :uri             (str "http://localhost:8080/documents/" id) 
                  :format          (ajax/json-request-format)
                  :response-format (ajax/json-response-format {:keywords? true}) 
                  :on-success      [::process-get-doc-by-id]
                  :on-failure      [::bad-response]
                  }
     :db  (assoc db :loading? true)}))

(re-frame/reg-event-fx
    ::upload-book
  (fn
    [{db :db} [_ form]]
    {:http-xhrio {:method          :post
                  :uri             "http://localhost:8080/books/"
                  :body form
                  :response-format (ajax/json-response-format {:keywords? true}) 
                  :on-success      [::get-documents]
                  :on-failure      [::bad-response]
                  }
     :db  (assoc db :loading? true)}))

(re-frame/reg-event-fx
    ::delete-document-by-id
  (fn
    [{db :db} [_ id]]
    {:http-xhrio {:method          :delete
                  :uri             (str "http://localhost:8080/documents/" id) 
                  :format          (ajax/json-request-format)
                  :response-format (ajax/json-response-format {:keywords? true}) 
                  :on-success      [::get-documents]
                  :on-failure      [::bad-response]
                  }
     :db  (assoc db :loading? true)}))
