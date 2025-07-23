import 'package:flutter/material.dart';

class TextStyleHelper {
  // Create text style with outline and shadow
  static TextStyle withOutlineAndShadow({
    required double fontSize,
    required FontWeight fontWeight,
    required Color color,
    Color outlineColor = Colors.black,
    double outlineWidth = 0.5,
    Color shadowColor = Colors.black26,
    double shadowBlurRadius = 2.0,
    Offset shadowOffset = const Offset(1, 1),
  }) {
    return TextStyle(
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color,
      shadows: [
        Shadow(
          color: shadowColor,
          blurRadius: shadowBlurRadius,
          offset: shadowOffset,
        ),
        // Add outline effect using multiple shadows
        Shadow(
          color: outlineColor,
          blurRadius: outlineWidth,
          offset: const Offset(-0.5, -0.5),
        ),
        Shadow(
          color: outlineColor,
          blurRadius: outlineWidth,
          offset: const Offset(0.5, -0.5),
        ),
        Shadow(
          color: outlineColor,
          blurRadius: outlineWidth,
          offset: const Offset(-0.5, 0.5),
        ),
        Shadow(
          color: outlineColor,
          blurRadius: outlineWidth,
          offset: const Offset(0.5, 0.5),
        ),
      ],
    );
  }

  // Predefined styles for different text types
  static TextStyle titleStyle({required Color color}) {
    return withOutlineAndShadow(
      fontSize: 20,
      fontWeight: FontWeight.bold,
      color: color,
      outlineWidth: 0.8,
      shadowBlurRadius: 3.0,
    );
  }

  static TextStyle valueStyle({required Color color}) {
    return withOutlineAndShadow(
      fontSize: 16,
      fontWeight: FontWeight.bold,
      color: color,
      outlineWidth: 0.6,
      shadowBlurRadius: 2.5,
    );
  }

  static TextStyle labelStyle({required Color color}) {
    return withOutlineAndShadow(
      fontSize: 14,
      fontWeight: FontWeight.w600,
      color: color,
      outlineWidth: 0.4,
      shadowBlurRadius: 2.0,
    );
  }

  static TextStyle pnlStyle({required Color color}) {
    return withOutlineAndShadow(
      fontSize: 16,
      fontWeight: FontWeight.bold,
      color: color,
      outlineWidth: 0.7,
      shadowBlurRadius: 2.5,
    );
  }

  static TextStyle exchangeStyle({required Color color}) {
    return withOutlineAndShadow(
      fontSize: 18,
      fontWeight: FontWeight.bold,
      color: color,
      outlineWidth: 0.6,
      shadowBlurRadius: 2.5,
    );
  }
}
