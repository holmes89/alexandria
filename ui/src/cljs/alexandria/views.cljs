(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]))


;; home

(defn home-panel []
  (let [name (re-frame/subscribe [::subs/name])]
    [:div
     [:h1.main-title "Alexandria" ]

     [:div
      [:a {:href "#/documents"}
       "documents"]]
     ]))


;; read

(defn read-panel []
  (let [doc (re-frame/subscribe [::subs/active-doc])]
    (fn []
      [:div.container
       [:h1.read-title (:display_name @doc)]])))

;; docs

(defn doc-icon
  [type]
  (if (= type "book")
    [:i.fas.fa-book]
    [:i.fas.fa-file-alt]))

(defn doc-item
  [{:keys [id display_name type]}]
  [:a.panel-block {:href (str "#/documents/" id)}
   [:div.doc-info [doc-icon type]
    display_name] ])


(defn doc-list []
  (let [docs (re-frame/subscribe [::subs/docs])]
    (fn []
      [:div.columns.is-mobile
       [:div.column.is-6.is-offset-3
        [:nav.panel
         [:p.panel-heading "Documents"]
         (for [doc @docs]
           ^{:key (:id doc)}[doc-item doc])]]])))


(defn doc-panel []
  (fn []
    [:div.container
     [:h1.main-title "Alexandria"]
     [doc-list]]))

;; main

(defn- panels [panel-name]
  (case panel-name
    :home-panel [home-panel]
    :doc-panel [doc-panel]
    :read-panel [read-panel]
    [:div]))

(defn show-panel [panel-name]
  [panels panel-name])

(defn main-panel []
  (let [active-panel (re-frame/subscribe [::subs/active-panel])]
    [show-panel @active-panel]))
