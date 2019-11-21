(ns alexandria.views
  (:require
   [re-frame.core :as re-frame]
   [alexandria.subs :as subs]
   [alexandria.events :as events]
   [react-pdf :as pdf]
   [clojure.string :as str]))


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
(defn file-upload-name []
  (peek (str/split (.-value (.getElementById js/document "file-upload")) "\\")))
(defn upload-form [] (js/FormData. (.getElementById js/document "upload")))

(defn upload-modal []
  (let [show (re-frame/subscribe [::subs/is-upload-showing?])
        upload-name (re-frame/subscribe [::subs/upload-file-name])]
    (if @show
      [:div.modal.is-active
       [:div.modal-background]
       [:div.modal-card
        [:header.modal-card-head
         [:p.modal-card-title "Upload"]
         [:button.delete {:on-click #(re-frame/dispatch [::events/hide-upload-modal])}]]
        [:section.modal-card-body
         [:form#upload
          [:div.field
           [:label.label "Name"]
           [:div.control
            [:input.input {:type "text" :name "name"}]]]
          [:div.field
           [:label.file-label
            [:input#file-upload.file-input {:type "file" :name "file" :on-change #(re-frame/dispatch [::events/update-upload-file-name (file-upload-name)])}]
            [:span.file-cta
             [:span.file-icon
              [:i.fas.fa-upload]]
             [:span.file-label "Choose a file..."]]
            [:span.file-name @upload-name]]]]]
        [:footer.modal-card-foot
         [:button.button.is-success {:on-click #((re-frame/dispatch [::events/upload-book (upload-form)])
                                                 (re-frame/dispatch [::events/hide-upload-modal]))}"Submit"]
         [:button.button {:on-click #(re-frame/dispatch [::events/hide-upload-modal])} "Cancel"]]]]
      [:div])))

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
     (navbar)
     [:a {:on-click #(re-frame/dispatch [::events/show-upload-modal])}
      [:i.fas.fa-plus] "Add"]
     [doc-list]
     (upload-modal)]))

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
