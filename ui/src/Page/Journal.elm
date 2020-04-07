module Page.Journal exposing (view)

import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)


view : Html msg
view =
    section [ class "hero is-large" ]
        [ div [ class "hero-body" ]
            [ div [ class "container" ] [ text "here" ]
            ]
        , div [ class "hero-foot" ]
            [ div [ class "container" ] [ text "here" ]
            ]
        ]
