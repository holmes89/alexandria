(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]
   [react-pdf :as pdf]))


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
    [:> pdf/Page {:pageNumber @page-num :scale @zoom}]))


(defn zoom-in []
  [:button.button {:on-click #(re-frame/dispatch [::events/zoom-in])}
   [:i.fas.fa-search-plus]])

(defn zoom-out []
  [:button.button {:on-click #(re-frame/dispatch [::events/zoom-out])}
   [:i.fas.fa-search-minus]])

(defn next-page []
  [:button.button {:on-click #(re-frame/dispatch [::events/next-page])}
   [:span "Next"]
   [:i.fas.fa-arrow-right]])

(defn prev-page []
  [:button.button {:on-click #(re-frame/dispatch [::events/prev-page])}
   [:i.fas.fa-arrow-left]
   [:span "Prev"]])

(defn pdf-reader [src]
  [:> pdf/Document {:file src}
   (pdf-page 1)])

(defn read-panel []
  (let [doc (re-frame/subscribe [::subs/active-doc])]
    (fn []
      [:div.container
       [:h1.read-title (:display_name @doc)]
       (prev-page)
       (next-page)
       (zoom-in)
       (zoom-out)
       (let [src (:path @doc)]
         (if src
           [pdf-reader (:path @doc)]))])))

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
