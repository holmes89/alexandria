(ns alexandria.db
  (:require
   [re-frame.core :as re-frame]))

;; -- Local Storage  ----------------------------------------------------------
;;
;; Part of the conduit challenge is to store a user in localStorage, and
;; on app startup, reload the user from when the program was last run.
;;
(def alex-user-key "alexandria-user")  ;; localstore key

(defn set-user-ls
  [user]
  (.setItem js/localStorage alex-user-key (str user)))  ;; sorted-map written as an EDN map

;; Removes user information from localStorge when a user logs out.
;;
(defn remove-user-ls
  []
  (.removeItem js/localStorage alex-user-key))

;; -- cofx Registrations  -----------------------------------------------------
;;
;; Use `reg-cofx` to register a "coeffect handler" which will inject the user
;; stored in localStorge.
;;
;; To see it used, look in `events.cljs` at the event handler for `:initialise-db`.
;; That event handler has the interceptor `(inject-cofx :local-store-user)`
;; The function registered below will be used to fulfill that request.
;;
;; We must supply a `sorted-map` but in localStorage it is stored as a `map`.
;;
(re-frame/reg-cofx
    ::local-store-user
  (fn [cofx _]
    (assoc cofx ::local-store-user  ;; put the local-store user into the coeffect under :local-store-user
           (into (sorted-map)      ;; read in user from localstore, and process into a sorted map
                 (some->> (.getItem js/localStorage alex-user-key)
                          (cljs.reader/read-string))))))  ;; EDN map -> map

(def default-db
  {:name "re-frame"
   :document-data []
   :active-doc "test"
   :page-num 1
   :zoom 1.0
   :show-upload false
   :upload-file-name ""
   :authenticated false})
