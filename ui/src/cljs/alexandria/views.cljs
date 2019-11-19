(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]
   [react-pdf :as pdf]))


;; shared components

(defn navbar []
  [:nav.navbar {:role "navigation" :aria-label "main navigation"}
   [:div.navbar-brand
    [:a.navbar-item {:href "#/documents"}
     [:span "Alexandria"]]]])

;; home

(defn home-panel []
  (let [name (re-frame/subscribe [::subs/name])]
    [:div
     [:h1.main-title "Alexandria" ]
     [:div
      [:a {:href "#/documents"}
       "documents"]]]))


;; read
(defn pdf-page [num]
  (let [page-num (re-frame/subscribe [::subs/doc-page])
        zoom (re-frame/subscribe [::subs/doc-zoom])]
    [:> pdf/Page {:pageNumber @page-num :scale @zoom :renderAnnotationLayer false}]))

(defn zoom-in []
  [:a {:on-click #(re-frame/dispatch [::events/zoom-in])}
   [:i.fas.fa-search-plus]])

(defn zoom-out []
  [:a { :on-click #(re-frame/dispatch [::events/zoom-out])}
   [:i.fas.fa-search-minus]])

(defn read-bar [title]
  [:nav#read-bar.navbar.is-dark
   [:div.navbar-menu
    [:div.navbar-start
     [:div.navbar-item
      [:h3.book-title title]]]
    [:div.navbar-end
     (zoom-in)
     (zoom-out)]]])

(defn next-page []
  [:a.page-turn {:on-click #(re-frame/dispatch [::events/next-page])}
   [:i.fas.fa-arrow-right]])

(defn prev-page []
  [:a.page-turn {:on-click #(re-frame/dispatch [::events/prev-page])}
   [:i.fas.fa-arrow-left]])

(defn pdf-reader [src]
  [:> pdf/Document {:file src}
   (pdf-page 1)])

(defn read-section []
  (let [doc (re-frame/subscribe [::subs/active-doc])]
    [:div
     (read-bar (:display_name @doc))
     [:div
      [:div.columns.is-gapless
       [:div.column.is-1
        (prev-page)]
       [:div#doc.column.is-10
        (let [src (:path @doc)]
          (if src
            [pdf-reader (:path @doc)]))]
       [:div.column.is-1
        (next-page)]]]]))


(defn read-panel []
  [:div
   (navbar)
   (read-section)])

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
