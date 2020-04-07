module Page.Journal exposing (view)

import Html exposing (..)
import Html.Attributes exposing (class, href, rows, src)


view : Html msg
view =
    section [ class "hero is-large" ]
        [ div [ class "hero-body" ]
            [ div [ class "container" ] [ text "here" ]
            ]
        , div [ class "hero-foot journal-draft-area" ]
            [ div [ class "container" ]
                [ div [ class "columns" ]
                    [ div [ class "column is-11" ]
                        [ textarea [ class "textarea", rows 5 ] [] ]
                    , div [ class "column" ]
                        [ button [ class "button" ] [ text "Submit" ] ]
                    ]
                ]
            ]
        ]
