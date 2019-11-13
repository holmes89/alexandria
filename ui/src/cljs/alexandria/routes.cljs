(ns alexandria.routes
  (:require-macros [secretary.core :refer [defroute]])
  (:import [goog History]
           [goog.history EventType])
  (:require
   [secretary.core :as secretary]
   [goog.events :as gevents]
   [re-frame.core :as re-frame]
   [alexandria.events :as events]))

(defn hook-browser-navigation! []
  (doto (History.)
    (gevents/listen
     EventType/NAVIGATE
     (fn [event]
       (secretary/dispatch! (.-token event))))
    (.setEnabled true)))

(defn app-routes []
  (secretary/set-config! :prefix "#")
  ;; --------------------
  ;; define routes here
  (defroute "/" []
    (re-frame/dispatch [::events/set-active-panel :home-panel]))

  (defroute "/documents/:id" [id]
    (re-frame/dispatch [::events/get-document-by-id id])
    (re-frame/dispatch [::events/set-active-panel :read-panel]))
  (defroute "/documents" []
    (re-frame/dispatch [::events/get-documents])
    (re-frame/dispatch [::events/set-active-panel :doc-panel]))


  ;; --------------------
  (hook-browser-navigation!))
