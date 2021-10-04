/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.spark.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterProjectedSources;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * V1alpha1SparkClusterSpecMasterProjected
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecMasterProjected {
  public static final String SERIALIZED_NAME_DEFAULT_MODE = "defaultMode";
  @SerializedName(SERIALIZED_NAME_DEFAULT_MODE)
  private Integer defaultMode;

  public static final String SERIALIZED_NAME_SOURCES = "sources";
  @SerializedName(SERIALIZED_NAME_SOURCES)
  private List<V1alpha1SparkClusterSpecMasterProjectedSources> sources = null;


  public V1alpha1SparkClusterSpecMasterProjected defaultMode(Integer defaultMode) {
    
    this.defaultMode = defaultMode;
    return this;
  }

   /**
   * Get defaultMode
   * @return defaultMode
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getDefaultMode() {
    return defaultMode;
  }


  public void setDefaultMode(Integer defaultMode) {
    this.defaultMode = defaultMode;
  }


  public V1alpha1SparkClusterSpecMasterProjected sources(List<V1alpha1SparkClusterSpecMasterProjectedSources> sources) {
    
    this.sources = sources;
    return this;
  }

  public V1alpha1SparkClusterSpecMasterProjected addSourcesItem(V1alpha1SparkClusterSpecMasterProjectedSources sourcesItem) {
    if (this.sources == null) {
      this.sources = new ArrayList<>();
    }
    this.sources.add(sourcesItem);
    return this;
  }

   /**
   * Get sources
   * @return sources
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<V1alpha1SparkClusterSpecMasterProjectedSources> getSources() {
    return sources;
  }


  public void setSources(List<V1alpha1SparkClusterSpecMasterProjectedSources> sources) {
    this.sources = sources;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecMasterProjected v1alpha1SparkClusterSpecMasterProjected = (V1alpha1SparkClusterSpecMasterProjected) o;
    return Objects.equals(this.defaultMode, v1alpha1SparkClusterSpecMasterProjected.defaultMode) &&
        Objects.equals(this.sources, v1alpha1SparkClusterSpecMasterProjected.sources);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultMode, sources);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecMasterProjected {\n");
    sb.append("    defaultMode: ").append(toIndentedString(defaultMode)).append("\n");
    sb.append("    sources: ").append(toIndentedString(sources)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}
