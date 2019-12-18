(ns alexandria.auth0
  (:require [re-frame.core :as re-frame]
            [alexandria.config :as config]
            [auth0-lock :as auth0]
            [alexandria.db :as db]))
(def lock
  "The auth0 lock instance used to login and make requests to Auth0"
  (let [client-id (:client-id config/auth0)
        domain (:domain config/auth0)
        options (clj->js {:auth {:responseType "token id_token"}})]
    (auth0/Auth0Lock. client-id domain options)))

(defn handle-profile-response [error profile] *
  "Handle the response for Auth0 profile request"
  (let [profile-clj (js->clj profile :keywordize-keys true)]
    (re-frame/dispatch [::set-user-profile profile-clj])))

(defn on-authenticated
  "Function called by auth0 lock on authentication"
  [auth-result-js]
  (let [auth-result-clj (js->clj auth-result-js :keywordize-keys true)
        access-token (:accessToken auth-result-clj)]
    (re-frame/dispatch [::set-auth-result auth-result-clj])
    (re-frame/dispatch [::set-authenticated true])
    (.getUserInfo lock access-token handle-profile-response)))

(.on lock "authenticated" on-authenticated)


;;; events

;; -- Interceptors --------------------------------------------------------------
;; Every event handler can be "wrapped" in a chain of interceptors. Each of these
;; interceptors can do things "before" and/or "after" the event handler is executed.
;; They are like the "middleware" of web servers, wrapping around the "handler".
;; Interceptors are a useful way of factoring out commonality (across event
;; handlers) and looking after cross-cutting concerns like logging or validation.
;;
;; They are also used to "inject" values into the `coeffects` parameter of
;; an event handler, when that handler needs access to certain resources.
;;
;; Each event handler can have its own chain of interceptors. Below we create
;; the interceptor chain shared by all event handlers which manipulate user.
;; A chain of interceptors is a vector.
;; Explanation of `trim-v` is given further below.
;;
(def set-user-interceptor [(re-frame/path :user)        ;; `:user` path within `db`, rather than the full `db`.
                           (re-frame/after db/set-user-ls) ;; write user to localstore (after)
                           ])            ;; removes first (event id) element from the event vec

;; After logging out clean up local-storage so that when a users refreshes
;; the browser she/he is not automatically logged-in, and because it's a
;; good practice to clean-up after yourself.
;;
(def remove-user-interceptor [(re-frame/after db/remove-user-ls)])


(re-frame/reg-event-db
    ::set-auth-result
  set-user-interceptor
  (fn [db [_ auth-result]]
    (assoc-in db [:user :auth-result] auth-result)))

(re-frame/reg-event-db
    ::set-authenticated
  (fn [db [_ authed]]
    (assoc db :authenticated authed)))

(re-frame/reg-event-db
    ::set-user-profile
  set-user-interceptor
  (fn [db [_ profile]]
    (assoc-in db [:user :profile] profile)))

(re-frame/reg-event-db
    ::logout
  remove-user-interceptor
  (fn [db [_ profile]]
    (dissoc db :user)
    (assoc db :authenticated false)))

;;; subscriptions
(re-frame/reg-sub
    ::user-name
  (fn [db]
    (get-in db [:user :profile :name])))

(re-frame/reg-sub
    ::token
  (fn [db]
    (get-in db [:user :auth-result :idToken])))

(re-frame/reg-sub
    ::authenticated
  (fn [db]
    (:authenticated db)))

(re-frame/reg-sub
    ::profile-image
  (fn [db]
    (get-in db [:user :profile :picture])))
