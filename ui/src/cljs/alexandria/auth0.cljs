(ns alexandria.auth0
  (:require [re-frame.core :as re-frame]
            [alexandria.config :as config]
            [auth0-lock :as auth0]))
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
(re-frame/reg-event-db
    ::set-auth-result
  (fn [db [_ auth-result]]
    (assoc-in db [:user :auth-result] auth-result)))

(re-frame/reg-event-db
    ::set-authenticated
  (fn [db [_ authed]]
    (assoc db :authenticated authed)))

(re-frame/reg-event-db
    ::set-user-profile
  (fn [db [_ profile]]
    (assoc-in db [:user :profile] profile)))

(re-frame/reg-event-db
    ::logout
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
