(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]))


;; home

(defn home-panel []
  (let [name (re-frame/subscribe [::subs/name])]
    [:div
     [:h1 (str "Hello from " @name ". This is the Home Page.")]

     [:div
      [:a {:href "#/about"}
       "go to About Page"]]
     ]))


;; about
(defn about-panel []
  [:div
   [:h1 "This foo is the About Page."]

   [:div
    [:a {:href "#/"}
     "go to Home Page"]]])

;; books

(defn book-item
  [{:keys [id display_name]}]
  [:a.panel-block {:href (str "#/books/" id)} display_name])


(defn book-list []
  (let [books (re-frame/subscribe [::subs/books])]
    (fn []
      [:div.columns.is-mobile
       [:div.column
        [:nav.panel
         [:p.panel-heading "Documents"]
         (for [book @books]
           ^{:key (:id book)}[book-item book])]]])))


(defn book-panel []
  (fn []
    [:div.container
     [:h1.main-title "Alexandria"]
     [book-list]]))

;; main

(defn- panels [panel-name]
  (case panel-name
    :home-panel [home-panel]
    :book-panel [book-panel]
    [:div]))

(defn show-panel [panel-name]
  [panels panel-name])

(defn main-panel []
  (let [active-panel (re-frame/subscribe [::subs/active-panel])]
    [show-panel @active-panel]))
