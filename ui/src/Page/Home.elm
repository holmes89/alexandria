module Page.Home exposing (view)

import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)


view : Html msg
view =
    section [ class "hero" ]
        [ div [ class "hero-body" ]
            [ div [ class "container" ]
                [ h1 [ class "title" ] [ text "Welcome!" ]
                ]
            ]
        ]
