(ns alexandria.pdfviewer
  (:require
   [reagent.core :as r]
   [pdfjs-dist :as pdfjs]))

;; empirically determined magic factor to remove white space below pdf
(def height-factor 1.13)

                                        ; PDFjs helper functions
                                        ; TODO: Do the grunt work only once instead of on every render
(defn render-page []
  (let [pdf (:pdf @app-state)
        current_page (get-in @app-state [:navigation :current_page 0])]
    (.then (.getPage pdf current_page)
           (fn [page]
             (let [desiredWidth (:pdf_width @app-state)
                   viewport (.getViewport page 1)
                   height (.-height viewport)
                   width (.-width viewport)
                   scale (/ desiredWidth width)
                   scaledViewport (.getViewport page scale)
                   canvas (-> js/document
                              (.querySelector "pdf-viewer")
                              (.querySelector "canvas"))
                   context (.getContext canvas "2d")
                   renderContext (js-obj "canvasContext" context "viewport" scaledViewport)]
                                        ; TODO: Eval employing the renderTask promise of PDFjs
               (aset canvas "height" (/ height height-factor))
               (.render page renderContext))))))


                                        ; OM Component helper functions
(defn valid-page? [page-num]
  (let [page-count (get-in @app-state [:navigation :page_count 0])]
    (and (> page-num 0) (<= page-num page-count))))

(defn render-page-if-valid [cursor f]
  (let [current_page (get-in cursor [:current_page 0])]
    (if (valid-page? (f current_page))
      (do
        (om/transact! cursor [:current_page 0] f)
        (render-page)))))

                                        ; OM Components
(defn pdf-navigation-position [cursor owner]
  (reify
    om/IRender
    (render [this]
      (let [current_page (get-in cursor [:current_page 0])
            page_count (get-in cursor [:page_count 0])]
        (dom/span #js {:className "pageCount" }
                  (if (= 0 page_count)
                    "Loading"
                    (str current_page " of " page_count)))))))

(defn pdf-navigation-buttons [cursor owner]
  (reify
    om/IRender
    (render [this]
      (let [current_page (get-in cursor [:current_page 0])]
        (dom/span #js {:className "navButtons" }
                  (dom/button #js {:onClick (fn [e]
                                              (render-page-if-valid cursor dec))}
                              "<")
                  (dom/button #js {:onClick (fn [e]
                                              (render-page-if-valid cursor inc))}
                              ">"))))))

(defn pdf-navigation-view [cursor owner]
  (reify
    om/IRender
    (render [this]
      (dom/div #js {:className "navigation"}
               (om/build pdf-navigation-buttons cursor)
               (om/build pdf-navigation-position cursor)))))

(defn pdfjs-viewer [cursor owner]
  (reify
    om/IDidMount
    (did-mount [_]
      (if (not (nil? (:pdf_workerSrc @app-state)))
        (do
          (aset js/PDFJS "workerSrc" (:pdf_workerSrc @app-state))))
      (let [loadingTask (.getDocument js/PDFJS (:pdf_url @app-state))]

        (aset loadingTask "onProgress" (fn [progress] (om/update! cursor
                                                                  [:progress :loading 0]
                                                                  (/ (.-loaded progress) (.-total progress)))))
        (.then loadingTask (fn [pdf]
                             (swap! app-state assoc :pdf pdf)
                             (swap! app-state update-in [:navigation :page_count 0] #(.-numPages pdf))
                             (render-page)))))
    om/IRender
    (render [this]
      (dom/canvas #js {:width (:pdf_width @app-state)}))))

(defn pdf-progress-view [cursor owner]
  (reify
    om/IRender
    (render [this]
      (let [progress (get-in cursor [:loading 0])]
        (dom/span #js {:className "progress" }
                  (if (not (= progress 1))
                    (str (js/parseInt (* 100 progress)) "%")))))))

(defn pdf-component-view [cursor owner]
  (reify
    om/IRender
    (render [this]
      (dom/div #js {:id "om-root"
                    :style #js { :width (:pdf_width @app-state)}}
               (dom/div #js {:className "menu" }
                        (om/build pdf-navigation-view (cursor :navigation)))
               (om/build pdfjs-viewer cursor)
               (om/build pdf-progress-view (cursor :progress))))))
